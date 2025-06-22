package dbx

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/type/decimal"
)

type entitiesStruct struct {
	EntityTypeCache        sync.Map
	TableNameCache         sync.Map
	cacheCreateEntityType  sync.Map
	cacheTableNameEntity   map[string]*EntityType
	hashCheckIsDbFieldAble map[reflect.Type]bool
}

// extract type from entity, entity can be a struct or a pointer to a struct or even an instance of a struct
func (e *entitiesStruct) GetType(entity interface{}) (*reflect.Type, error) {
	if entity == nil {
		return nil, fmt.Errorf("entity is nil")
	}
	if rt, ok := entity.(reflect.Type); ok {
		if rt.Kind() == reflect.Ptr {
			rt = rt.Elem()
		}
		if rt.Kind() == reflect.Slice {
			rt = rt.Elem()
		}
		if rt.Kind() == reflect.Ptr {
			rt = rt.Elem()
		}
		if rt.Kind() != reflect.Struct {
			return nil, fmt.Errorf("entity must be a struct or a pointer to a struct")
		}
		return &rt, nil
	}
	typ := reflect.TypeOf(entity)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("entity must be a struct or a pointer to a struct")
	}
	return &typ, nil

}
func (e *entitiesStruct) GetEntityType(entity interface{}) (*EntityType, error) {
	typ, err := e.GetType(entity)
	if err != nil {
		return nil, err
	}
	return e.newEntityType(*typ)
}

// create or get (if exists)new EntityType from reflect.Type
func (e *entitiesStruct) newEntityType(t reflect.Type) (*EntityType, error) {
	key := t.PkgPath() + t.Name()
	if v, ok := e.EntityTypeCache.Load(key); ok {
		return v.(*EntityType), nil
	}
	ret, err := e.newEntityTypeNoCache(t)
	if err != nil {
		return nil, err
	}
	e.EntityTypeCache.Store(key, ret)
	return ret, nil
}
func (e *entitiesStruct) GetTableNameByType(t reflect.Type) string {
	key := t.PkgPath() + t.Name()
	if v, ok := e.TableNameCache.Load(key); ok {
		return v.(string)
	}
	ret := e.getTableNameNoCache(t)
	e.TableNameCache.Store(key, ret)
	return ret
}
func (e *entitiesStruct) getTableNameNoCache(t reflect.Type) string {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Anonymous {
			if f.Type == reflect.TypeOf(EntityModel{}) {
				return t.Name()
			} else {
				return e.GetTableNameByType(f.Type)

			}
		}
	}
	return ""

}
func (e *entitiesStruct) newEntityTypeNoCache(t reflect.Type) (*EntityType, error) {
	//check cache
	tableName := e.GetTableNameByType(t)
	if tableName == "" {
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		if t.Kind() == reflect.Slice {
			t = t.Elem()
		}
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		return nil, fmt.Errorf("The entity type %s is not a valid entity type, please embed EntityModel and  add `db:\"table_name\"` tag to the entity struct", t.PkgPath()+"."+t.Name())
	}

	ret := EntityType{
		Type:         t,
		TableName:    tableName,
		fieldMap:     sync.Map{},
		RefEntities:  []*EntityType{},
		EntityFields: []*EntityField{},
	}
	fields, refTable := e.GetAllFields(ret.Type)

	ret.EntityFields = make([]*EntityField, 0)
	for _, f := range fields {
		nf := f.Type
		if nf.Kind() == reflect.Ptr {
			nf = nf.Elem()
		}
		if nf.Kind() == reflect.Slice {
			nf = nf.Elem()
		}
		if nf.Kind() == reflect.Ptr {
			nf = nf.Elem()
		}

		ef := EntityField{
			TableName:       ret.TableName,
			StructField:     f,
			AllowNull:       true,
			NonPtrFieldType: nf,
		}
		err := ef.initPropertiesByTags()
		if err != nil {
			return nil, err
		}

		ret.EntityFields = append(ret.EntityFields, &ef)
	}
	for _, ref := range refTable {
		refType := ref.Type
		if refType.Kind() == reflect.Ptr {
			refType = refType.Elem()
		}
		if refType.Kind() == reflect.Slice {
			refType = refType.Elem()
		}
		if refType.Kind() == reflect.Ptr {
			refType = refType.Elem()
		}
		refEntity, err := e.newEntityType(refType)
		if err != nil {
			return nil, err
		}
		refEntityField := EntityField{
			StructField: ref,
		}
		err = refEntityField.initPropertiesByTags()

		fkNameList := strings.Split(refEntityField.ForeignKey, ",")
		name := refEntityField
		fmt.Println(name)
		for _, fkName := range fkNameList {
			fx := refEntity.GetFieldByName(fkName)
			if fx == nil {
				return nil, fmt.Errorf("invalid foreign key: %s.%s tag in models %s is %s", refEntity.TableName, fkName, t.Name(), ref.Tag)
			}
			refEntity.RefFields = append(refEntity.RefFields, fx)
		}

		if err != nil {
			return nil, err
		}
		ret.RefEntities = append(ret.RefEntities, refEntity)
	}

	return &ret, nil
}

// get enties by table name
func (e *entitiesStruct) GetEntityTypeByTableName(tableName string) *EntityType {
	if ret, ok := e.cacheTableNameEntity[tableName]; ok {
		return ret
	}
	return nil
}

// load all fields of the entity type, including embedded fields. all fields can be used for database operation.
// @return all fields, all reference fields, error
func (e *entitiesStruct) GetAllFields(typ reflect.Type) ([]reflect.StructField, []reflect.StructField) {
	ret := make([]reflect.StructField, 0)
	check := map[string]bool{}
	anonymousFields := []reflect.StructField{}
	refField := []reflect.StructField{}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Type == reflect.TypeOf(FullTextSearchColumn("")) { //loi  FullTextSearchColumn (type) is not an expressioncompiler
			ret = append(ret, field)
			continue
		}
		if field.Anonymous {

			anonymousFields = append(anonymousFields, field)

			continue
		} else {
			ft := field.Type

			if field.Type.Kind() == reflect.Ptr {
				ft = field.Type.Elem()
			}
			if ft.Kind() == reflect.Slice {
				ft = ft.Elem()
			}
			if field.Type.Kind() == reflect.Ptr {
				ft = field.Type.Elem()
			}
			if _, ok := e.hashCheckIsDbFieldAble[ft]; !ok {

				refField = append(refField, field)
				continue
			}
			check[field.Name] = true
			ret = append(ret, field)
		}
	}
	for _, field := range anonymousFields {
		fields, _ := e.GetAllFields(field.Type)

		for _, f := range fields {
			if _, ok := check[f.Name]; !ok { //check if field is not exist
				check[f.Name] = true
				ret = append(ret, f)
			}

		}
	}

	return ret, refField
}
func (e *entitiesStruct) CreateEntityType(entity interface{}) (*EntityType, error) {
	if entity == nil {
		return nil, fmt.Errorf("entity type must not be nil")
	}
	if ft, ok := entity.(reflect.Type); ok {
		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}
		if ft.Kind() == reflect.Slice { //in case of slice
			ft = ft.Elem()

		}
		if ft.Kind() == reflect.Ptr { // in case of slice of pointer
			ft = ft.Elem()
		}
		if ft.Kind() != reflect.Struct { //in case of slice of non-struct
			return nil, fmt.Errorf("entity type must be a struct or a slice of struct, but got %v", ft.Kind())
		}
		key := ft.PkgPath() + "-" + ft.Name()
		//check cache
		if retEntity, ok := e.cacheCreateEntityType.Load(key); ok {
			return retEntity.(*EntityType), nil
		}

		retEntity, err := e.newEntityType(ft)
		if err != nil {
			return nil, err
		}
		//save to cache
		e.cacheCreateEntityType.Store(ft, retEntity)

		return retEntity, nil
	}
	typ := reflect.TypeOf(entity)
	if typ.Kind() == reflect.Ptr { // in case of pointer
		typ = typ.Elem()
	}
	if typ.Kind() == reflect.Slice { //in case of slice
		typ = typ.Elem()

	}
	if typ.Kind() == reflect.Ptr { // in case of slice of pointer
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct { //in case of slice of non-struct
		return nil, fmt.Errorf("entity type must be a struct or a slice of struct, but got %v", typ.Kind())
	}
	key := typ.PkgPath() + "-" + typ.Name()
	//check cache
	if retEntity, ok := e.cacheCreateEntityType.Load(key); ok {
		return retEntity.(*EntityType), nil
	}
	ret, err := e.newEntityType(typ)
	if err != nil {
		return nil, err
	}
	e.cacheTableNameEntity[ret.TableName] = ret
	//save to cache
	e.cacheCreateEntityType.Store(key, ret)
	uk := map[string][]string{}
	for _, f := range ret.EntityFields {
		if f.UkName != "" {
			if _, ok := uk[f.UkName]; !ok {
				uk[f.UkName] = []string{f.Name}
			} else {
				uk[f.UkName] = append(uk[f.UkName], f.Name)
			}

		}
	}
	for k, v := range uk {
		dbxEntityCache.set_uk(ret.TableName+"_"+k, v)
	}

	return ret, nil
}

// Entities is a global variable of entitiesStruct a utility struct for managing entity types and their fields.
var Entities *entitiesStruct

func init() {
	Entities = &entitiesStruct{
		EntityTypeCache:       sync.Map{},
		TableNameCache:        sync.Map{},
		cacheCreateEntityType: sync.Map{},
		cacheTableNameEntity:  map[string]*EntityType{},
	}
	Entities.hashCheckIsDbFieldAble = map[reflect.Type]bool{
		reflect.TypeOf(int(0)):      true,
		reflect.TypeOf(int8(0)):     true,
		reflect.TypeOf(int16(0)):    true,
		reflect.TypeOf(int32(0)):    true,
		reflect.TypeOf(int64(0)):    true,
		reflect.TypeOf(uint(0)):     true,
		reflect.TypeOf(uint8(0)):    true,
		reflect.TypeOf(uint16(0)):   true,
		reflect.TypeOf(uint32(0)):   true,
		reflect.TypeOf(uint64(0)):   true,
		reflect.TypeOf(float32(0)):  true,
		reflect.TypeOf(float64(0)):  true,
		reflect.TypeOf(string("")):  true,
		reflect.TypeOf(bool(false)): true,
		reflect.TypeOf(time.Time{}): true,

		reflect.TypeOf(decimal.Decimal{}): true,
		reflect.TypeOf(uuid.UUID{}):       true,
	}

}
