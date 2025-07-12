package migrate

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type modelRegistryInfo struct {
	tableName string
	modelType reflect.Type
	entity    Entity
}
type modelRegister struct {
	cacheModelRegistry  sync.Map
	cacheGetModelByType sync.Map
}

func (reg *modelRegister) getTableName(typ reflect.Type) (string, error) {
	// scan field
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Type.Field(0).Type == reflect.TypeOf(Entity{}) {
			tagInfo := field.Tag.Get("db")
			if strings.HasPrefix(tagInfo, "table:") {
				return tagInfo[6:], nil
			}
			return utilsInstance.Pluralize(utilsInstance.SnakeCase(typ.Name())), nil

		}
	}
	return "", fmt.Errorf("model %s has no table tag", typ.String())
}
func (reg *modelRegister) Add(m ...interface{}) {
	for _, model := range m {

		typ := reflect.TypeOf(model)
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		tableName, err := reg.getTableName(typ)
		if err != nil {
			panic(err)
		}
		cols, err := utilsInstance.ParseStruct(typ)
		if err != nil {
			panic(err)
		}
		entity := Entity{
			entityType: typ,
			tableName:  tableName,
			cols:       cols,
		}

		entity.primaryConstraints = utilsInstance.GetPrimaryKey(&entity)
		entity.uniqueConstraints = make(map[string][]ColumnDef)
		entity.indexConstraints = make(map[string][]ColumnDef)

		cacheItem := modelRegistryInfo{
			tableName: tableName,
			modelType: typ,
			entity:    entity,
		}

		reg.cacheModelRegistry.Store(typ, cacheItem)
	}
}
func (reg *modelRegister) GetAllModels() []modelRegistryInfo {
	ret := make([]modelRegistryInfo, 0)
	reg.cacheModelRegistry.Range(func(key, value interface{}) bool {
		ret = append(ret, value.(modelRegistryInfo))
		return true
	})
	return ret
}
func (reg *modelRegister) GetModelByType(typ reflect.Type) *modelRegistryInfo {
	if v, ok := reg.cacheGetModelByType.Load(typ); ok {
		return v.(*modelRegistryInfo)
	}
	var ret *modelRegistryInfo
	reg.cacheModelRegistry.Range(func(key, value interface{}) bool {
		m := value.(modelRegistryInfo)
		ret = &m
		return true
	})
	reg.cacheGetModelByType.Store(typ, ret)
	return ret
}

var ModelRegistry = &modelRegister{}
