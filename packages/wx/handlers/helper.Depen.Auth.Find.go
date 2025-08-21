package handlers

import (
	"reflect"
)

var IsGeneriAuth func(typ reflect.Type) bool

func (h *helperType) DepenAuthFind(typ reflect.Type) ([]int, error) {
	return h.depenFindAuthWithVisited(typ, map[reflect.Type]bool{})
}
func (h *helperType) depenFindAuthWithVisited(typ reflect.Type, visisted map[reflect.Type]bool) ([]int, error) {

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil, nil
	}
	if _, ok := visisted[typ]; ok {
		return nil, nil
	}
	visisted[typ] = true
	if IsGeneriAuth(typ) {
		return []int{}, nil
	}
	for i := 0; i < typ.NumField(); i++ {
		fieldTYpe := typ.Field(i).Type
		indexes, err := h.depenFindAuthWithVisited(fieldTYpe, visisted)
		if err != nil {
			return nil, err
		}
		if indexes != nil {
			return append(typ.Field(i).Index, indexes...), nil
		}

	}
	return nil, nil
}
