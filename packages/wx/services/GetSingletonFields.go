package services

import (
	"reflect"
	"sync"
)

type initGetSingletonFields struct {
	once     sync.Once
	instance [][]int
}

var cacheGetSingletonFields = sync.Map{}

func (svc *serviceUtilsType) GetSingletonFields(typ reflect.Type) [][]int {
	actual, _ := cacheGetSingletonFields.LoadOrStore(typ, &initGetSingletonFields{})
	initService := actual.(*initGetSingletonFields)
	initService.once.Do(func() {
		initService.instance = svc.getSingletonFieldsInternal(typ, map[reflect.Type]bool{})
	})
	return initService.instance
}

/*
This function will detect all fields of type if Singleton type return list of field index
*/
func (svc *serviceUtilsType) getSingletonFieldsInternal(typ reflect.Type, visitedType map[reflect.Type]bool) [][]int {
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

		if svc.IsFieldSingleton(field) {
			ret = append(ret, field.Index)

		} else {
			fieldType := field.Type

			subRet := svc.getSingletonFieldsInternal(fieldType, visitedType)
			if len(subRet) > 0 {
				for _, x := range subRet {
					ret = append(ret, append(field.Index, x...))
				}
			}
		}
	}
	return ret

}

type initCreateSingletonInstance struct {
	once          sync.Once
	valueInstance reflect.Value
}

var cacheCreateSingletonInstance = sync.Map{}

func (svc *serviceUtilsType) CreateSingletonInstance(typ reflect.Type) reflect.Value {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	actual, _ := cacheCreateSingletonInstance.LoadOrStore(typ, &initCreateSingletonInstance{})
	initService := actual.(*initCreateSingletonInstance)
	initService.once.Do(func() {
		initService.valueInstance = reflect.New(typ)
	})
	return initService.valueInstance
}
