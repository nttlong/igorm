package handlers

import (
	"reflect"
	"sync"
)

var IsGenericDepen func(typ reflect.Type) bool

type initIsGenericDepen struct {
	IsGeneric bool
	once      sync.Once
}

var cacheIsGenericDepen sync.Map

func (h *helperType) IsGenericDepen(typ reflect.Type) bool {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	actual, _ := cacheIsGenericDepen.LoadOrStore(typ, &initIsGenericDepen{})
	init := actual.(*initIsGenericDepen)
	init.once.Do(func() {
		init.IsGeneric = IsGenericDepen(typ)
	})

	return init.IsGeneric
}

type initGetDepen struct {
	Depen [][]int
	once  sync.Once
}

var cacheGetDepen sync.Map

func (h *helperType) GetDepen(typ reflect.Type) [][]int {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	actual, _ := cacheGetDepen.LoadOrStore(typ, &initGetDepen{})
	init := actual.(*initGetDepen)
	init.once.Do(func() {
		init.Depen = h.getDepen(typ)
	})

	return init.Depen
}
func (h *helperType) getDepen(typ reflect.Type) [][]int {
	ret := [][]int{}
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil
	}
	if h.IsGenericDepen(typ) {
		return ret
	}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		indexes := h.GetDepen(field.Type)
		if indexes != nil {
			if len(indexes) > 0 {
				for _, index := range indexes {
					fieldIndex := append(field.Index, index...)
					ret = append(ret, fieldIndex)
				}
			} else {
				ret = append(ret, field.Index)
			}
		}
	}
	return ret
}
