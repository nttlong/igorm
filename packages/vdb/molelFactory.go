package vdb

import "reflect"

type modelFacoryType struct {
}

func (m *modelFacoryType) CreateFromType(typ reflect.Type) (interface{}, error) {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return reflect.New(typ).Elem().Interface(), nil
}

var modelFacory = &modelFacoryType{}
