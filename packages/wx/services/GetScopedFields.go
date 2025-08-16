package services

import (
	"reflect"
	"sync"
)

type initGetScopeFields struct {
	once     sync.Once
	instance [][]int
}

var cacheGetScopeFields = sync.Map{}

func (svc *serviceUtilsType) GetGetScopeFields(typ reflect.Type) [][]int {
	actual, _ := cacheGetScopeFields.LoadOrStore(typ, &initGetScopeFields{})
	initService := actual.(*initGetScopeFields)
	initService.once.Do(func() {
		initService.instance = svc.getScopedFieldsInternal(typ, map[reflect.Type]bool{})
	})
	return initService.instance
}

/*
This function will detect all fields of type if Singleton type return list of field index
*/
func (svc *serviceUtilsType) getScopedFieldsInternal(typ reflect.Type, visitedType map[reflect.Type]bool) [][]int {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil
	}
	if visitedType[typ] {
		return nil
	}
	visitedType[typ] = true

	ret := [][]int{}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		if svc.IsFieldScoped(field) {
			ret = append(ret, field.Index)

		} else {
			fieldType := field.Type

			subRet := svc.getScopedFieldsInternal(fieldType, visitedType)
			if len(subRet) > 0 {
				for _, x := range subRet {
					ret = append(ret, append(field.Index, x...))
				}
			}
		}
	}
	return ret

}
