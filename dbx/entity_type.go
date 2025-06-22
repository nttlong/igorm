package dbx

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

type EntityType struct {
	reflect.Type
	TableName                       string
	fieldMap                        sync.Map
	RefEntities                     []*EntityType
	EntityFields                    []*EntityField
	IsLoaded                        bool
	RefFields                       []*EntityField
	mapCols                         map[string]*EntityField
	hasMapCols                      bool // map of field name to EntityField
	autoValueColsName               []string
	hasGenerateDefaultValueColsName bool // list of field names that have auto value
	hasGetPrimaryKeyName            bool // list of field names that have primary key
	defaultValueColsNames           []string
	primaryKeyNames                 []string
	hasGetDefaultValueColsName      bool // list of field names that have default value
	pkAutoCols                      []string
	hasPkAutoCols                   bool // list of field names that have primary key and auto value
}
type EntityField struct {
	TableName string
	reflect.StructField
	AllowNull    bool
	IsPrimaryKey bool

	DefaultValue    string
	MaxLen          int
	ForeignKey      string
	IndexName       string
	UkName          string
	NonPtrFieldType reflect.Type

	IsFullTextSearch bool
	IsBSON           bool
}

//var newEntityTypeCache = sync.Map{}

func (f *EntityField) initPropertiesByTags() error {
	if f.Type.Kind() == reflect.Ptr {
		f.AllowNull = true
	} else {
		f.AllowNull = false
	}

	strTags := ";" + f.Tag.Get("db") + ";"
	f.MaxLen = -1
	ft := f.Type
	if f.Type == reflect.TypeOf(FullTextSearchColumn("")) {

		f.IsFullTextSearch = true
		return nil
	}
	if f.Type.Kind() == reflect.Ptr {
		ft = f.Type.Elem()
	}
	f.NonPtrFieldType = ft
	for k, v := range replacerConstraint {
		for _, t := range v {

			strTags = strings.ReplaceAll(strTags, ";"+t+";", ";"+k+";")
			strTags = strings.ReplaceAll(strTags, ";"+t+":", ";"+k+":")
			strTags = strings.ReplaceAll(strTags, ";"+t+"(", ";"+k+"(")

		}

	}
	if f.Type.Kind() == reflect.Ptr {
		f.AllowNull = true
	}
	tags := strings.Split(strTags, ";")
	for _, tag := range tags {
		if tag == "" {
			continue
		}
		if tag == "pk" {
			f.IsPrimaryKey = true

		}
		if tag == "auto" {
			f.DefaultValue = "auto"
		}
		if strings.HasPrefix(tag, "size:") {
			strSize := tag[5:]
			intSize, err := strconv.Atoi(strSize)
			if err != nil {
				return fmt.Errorf("invalid size tag: %s", strTags)
			}
			f.MaxLen = intSize
		}
		if strings.HasPrefix(tag, "df:") {
			f.DefaultValue = tag[3:]
		}
		if strings.HasPrefix(tag, "fk:") {
			f.ForeignKey = tag[3:]

		}
		if strings.HasPrefix(tag, "fk(") && strings.HasSuffix(tag, ")") {
			f.ForeignKey = tag[3 : len(tag)-1]
		}
		if strings.HasPrefix(tag, "idx") {
			indexName := f.Name + "_idx"
			if strings.Contains(tag, ":") {
				indexName = tag[4:]

			}
			f.IndexName = indexName

		}
		if strings.HasPrefix(tag, "uk") {
			f.UkName = f.Name + "_uk"
			if strings.Contains(tag, ":") {
				f.UkName = tag[4:]
			}
		}
		if strings.HasPrefix(tag, "vachar(") && strings.HasSuffix(tag, ")") {
			strLen := tag[7 : len(tag)-1]
			intLen, err := strconv.Atoi(strLen)
			if err != nil {
				return fmt.Errorf("invalid vachar tag: %s", strTags)
			}
			f.MaxLen = intLen
		}
		if strings.HasPrefix(tag, "nvarchar(") && strings.HasSuffix(tag, ")") {
			strLen := tag[9 : len(tag)-1]
			intLen, err := strconv.Atoi(strLen)
			if err != nil {
				return fmt.Errorf("invalid nvarchar tag: %s", strTags)
			}
			f.MaxLen = intLen
		}
		if strings.HasPrefix(tag, "text(") && strings.HasSuffix(tag, ")") {
			strLen := tag[5 : len(tag)-1]
			intLen, err := strconv.Atoi(strLen)
			if err != nil {
				return fmt.Errorf("invalid nvarchar tag: %s", strTags)
			}
			f.MaxLen = intLen
		}

	}

	return nil

}

// func (e *EntityType) getAllFieldsDelete() ([]*EntityField, error) {
// 	// if e.IsLoaded {
// 	// 	return e.EntityFields, nil
// 	// }
// 	//check cache

// 	fields, refFields := getAllFields(e.Type)

// 	// sort fields by field index
// 	sort.Slice(fields, func(i, j int) bool {
// 		return fields[i].Index[0] < fields[j].Index[0]
// 	})
// 	ret := make([]*EntityField, 0)
// 	for _, field := range fields {

// 		ef := EntityField{
// 			StructField: field,
// 		}
// 		err := ef.initPropertiesByTags()
// 		if err != nil {
// 			return nil, err
// 		}
// 		ret = append(ret, &ef)
// 	}
// 	eRefEntity := make([]*EntityType, 0)
// 	for _, field := range refFields {
// 		ft := field.Type
// 		if ft.Kind() == reflect.Ptr {
// 			ft = ft.Elem()
// 		}
// 		if ft.Kind() == reflect.Slice {
// 			ft = ft.Elem()
// 		}
// 		if ft.Kind() == reflect.Ptr {
// 			ft = ft.Elem()
// 		}
// 		ef := EntityType{
// 			Type:        ft,
// 			TableName:   ft.Name(),
// 			fieldMap:    sync.Map{},
// 			RefEntities: []*EntityType{},
// 		}
// 		eRefEntity = append(eRefEntity, &ef)

// 	}
// 	e.RefEntities = eRefEntity
// 	//save to cache
// 	e.IsLoaded = true
// 	return ret, nil
// }

var (
	lockGetFieldByName sync.Map
)

func (e *EntityType) GetFieldByName(FieldName string) *EntityField {
	//check cache
	FieldName = strings.ToLower(FieldName)
	if field, ok := e.fieldMap.Load(FieldName); ok {
		ret := field.(EntityField)
		return &ret
	}

	for _, f := range e.EntityFields {
		if strings.EqualFold(f.Name, FieldName) {
			e.fieldMap.Store(FieldName, *f)
			return f
		}
	}
	return nil

}

func (e *EntityType) GetPrimaryKey() []*EntityField {

	ret := make([]*EntityField, 0)
	for _, field := range e.EntityFields {
		if field.IsPrimaryKey {
			ret = append(ret, field)
		}
	}
	return ret
}
func (e *EntityType) GetPkAutoCos() []string {
	if !e.hasPkAutoCols {
		ret := []string{}
		for _, field := range e.EntityFields {
			if field.IsPrimaryKey && field.DefaultValue == "auto" {
				ret = append(ret, field.Name)
			}
		}
		e.pkAutoCols = ret
		e.hasPkAutoCols = true
		return ret
	}
	return e.pkAutoCols

}
func (e *EntityType) GetPrimaryKeyName() []string {
	if !e.hasGetPrimaryKeyName {
		ret := []string{}
		for _, field := range e.EntityFields {
			if field.IsPrimaryKey {
				ret = append(ret, field.Name)
			}
		}
		e.primaryKeyNames = ret
		e.hasGetPrimaryKeyName = true
		return ret
	}
	return e.primaryKeyNames

}
func (e *EntityType) GetForeignKey() []*EntityField {

	ret := make([]*EntityField, 0)
	for _, field := range e.EntityFields {
		if field.ForeignKey != "" {
			ret = append(ret, field)
		}
	}
	return ret
}
func (e *EntityType) GetNonKeyFields() []EntityField {

	ret := make([]EntityField, 0)
	for _, field := range e.EntityFields {
		if !field.IsPrimaryKey {
			ret = append(ret, *field)
		}
	}
	return ret
}

// get index return map[indexName]EntityField
func (e *EntityType) GetIndex() map[string][]*EntityField {
	ret := map[string][]*EntityField{}

	for _, field := range e.EntityFields {
		if field.IndexName != "" && field.UkName == "" {
			//check if index already exist
			if fields, ok := ret[field.IndexName]; ok {
				fields = append(fields, field)
				ret[field.IndexName] = fields
			} else {
				ret[field.IndexName] = []*EntityField{field}
			}

		}
	}
	return ret
}
func (e *EntityType) GetUniqueKey() map[string][]*EntityField {
	ret := map[string][]*EntityField{}

	for _, field := range e.EntityFields {
		if field.UkName != "" {
			//check if index already exist
			if fields, ok := ret[field.UkName]; ok {
				fields = append(fields, field)
				ret[field.UkName] = fields
			} else {
				ret[field.UkName] = []*EntityField{field}
			}

		}
	}
	return ret
}

func (e *EntityType) getDefaultValueColsNames() []string {
	if !e.hasGenerateDefaultValueColsName {
		ret := []string{}
		for _, field := range e.EntityFields {
			if field.DefaultValue != "" && field.DefaultValue != "auto" && !field.IsPrimaryKey {
				ret = append(ret, field.Name)
			}
		}
		e.defaultValueColsNames = ret
		e.hasGenerateDefaultValueColsName = true

		return ret
	}
	return e.defaultValueColsNames
}
func (e *EntityType) getMapCols() map[string]*EntityField {
	if !e.hasMapCols {
		e.mapCols = make(map[string]*EntityField, 0)
		for _, field := range e.EntityFields {
			e.mapCols[field.Name] = field
		}
		e.hasMapCols = true
	}
	return e.mapCols
}

type ForeignKeyInfo struct {
	FromEntity *EntityType
	FromFields []*EntityField
	ToEntity   *EntityType
	ToFields   []*EntityField
}
type fkInfoEntry struct {
	ForeignTable  string
	OwnerTable    string
	ForeignFields []string
	OwnerFields   []string
}

func (e *EntityType) GetForeignKeyRef() map[string]fkInfoEntry {
	// retList := []*ForeignKeyInfo{}
	mapRefEntities := map[string]fkInfoEntry{}
	fmt.Println(e.TableName)

	for i := 0; i < e.Type.NumField(); i++ {
		eField := e.Type.Field(i)

		tag := ";" + eField.Tag.Get("db") + ";"
		// fmt.Println(tag)

		if strings.Contains(tag, ";fk:") || strings.Contains(tag, ";fk(") {
			fkField := ""
			if strings.Contains(tag, ";fk:") {
				fkField = strings.Split(tag, ";fk:")[1]
				fkField = strings.Split(fkField, ";")[0]
			}
			if strings.Contains(tag, ";fk(") {
				fkField = strings.Split(tag, ";fk(")[1]
				fkField = strings.Split(fkField, ")")[0]
			}
			// fmt.Println(fkField)
			if fkField != "" {
				eFieldTYpe := eField.Type
				if eFieldTYpe.Kind() == reflect.Ptr {
					eFieldTYpe = eFieldTYpe.Elem()
				}
				if eFieldTYpe.Kind() == reflect.Slice {
					eFieldTYpe = eFieldTYpe.Elem()
				}
				if eFieldTYpe.Kind() == reflect.Ptr {
					eFieldTYpe = eFieldTYpe.Elem()
				}
				for j := 0; j < eFieldTYpe.NumField(); j++ {
					refField := eFieldTYpe.Field(j)
					if refField.Anonymous && refField.Type == reflect.TypeOf(EntityModel{}) {
						fkTable := eFieldTYpe.Name()
						fkName := e.TableName + "_" + fkTable
						fromFields := e.GetPrimaryKey()
						if len(fromFields) > 1 {
							panic("not support multiple primary key")
						}
						if len(fromFields) == 0 {
							panic("entity must have a primary key")
						}
						if entry, ok := mapRefEntities[fkName]; ok {
							// entry.FromFields = append(entry.FromFields, fromFields[0].Name)
							// entry.ToFields = append(entry.ToFields, refField.Name)
							mapRefEntities[fkName] = entry
							// fmt.Println(entry)

						} else {
							mapRefEntities[fkName] = fkInfoEntry{
								ForeignTable:  e.TableName,
								OwnerTable:    fkTable,
								ForeignFields: []string{fromFields[0].Name},
								OwnerFields:   []string{fkField},
							}

						}
					}
				}

			}

		}

	}
	// for _, val := range mapRefEntities {
	// 	fmt.Println(val)
	// 	// refEntity := val
	// 	// ret := ForeignKeyInfo{}
	// 	// ret.FromEntity = refEntity
	// 	// ret.ToEntity = e
	// 	// ret.FromFields = refEntity.RefFields
	// 	// ret.ToFields = e.GetPrimaryKey()
	// 	// retList = append(retList, &ret)
	// }

	// for _, refEntity := range e.RefEntities {
	// 	ret := ForeignKeyInfo{}
	// 	ret.FromEntity = refEntity
	// 	ret.ToEntity = e
	// 	ret.FromFields = refEntity.RefFields
	// 	ret.ToFields = e.GetPrimaryKey()
	// 	retList = append(retList, &ret)

	// }
	return mapRefEntities

}

// Get all fields of the entity type, including embedded fields.

var replacerConstraint = map[string][]string{
	// "nvarchar": {"varchar"},
	"pk":   {"primary_key", "primarykey", "primary", "primary_key_constraint"},
	"fk":   {"foreign_key", "foreignkey", "foreign", "foreign_key_constraint"},
	"uk":   {"unique", "unique_key", "uniquekey", "unique_key_constraint"},
	"idx":  {"index", "index_key", "indexkey", "index_constraint"},
	"text": {"varchar", "varchar", "varchar2", "nvarchar"},

	"size": {"length", "len"},
	"df":   {"default", "default_value", "default_value_constraint"},
	"auto": {"auto_increment", "autoincrement", "serial_key", "serialkey", "serial_key_constraint"},
}
