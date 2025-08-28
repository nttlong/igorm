package vdb

import (
	"reflect"
)

type modelFacoryType struct {
}

func (m *modelFacoryType) CreateFromType(typ reflect.Type) (interface{}, error) {

	model := ModelRegistry.GetModelByType(typ)
	if model == nil {

		return nil, NewModelError(typ)
	}
	return reflect.New(typ).Elem().Interface(), nil
}

var modelFacory = &modelFacoryType{}
