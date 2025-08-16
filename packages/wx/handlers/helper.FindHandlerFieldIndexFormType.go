package handlers

import (
	"reflect"
	"wx/internal"
	"wx/services"
)

func (h *helperType) FindHandlerFieldIndexFormType(typ reflect.Type) ([]int, error) {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	key := typ.String() + "/helperType/FindHandlerFieldIndexFormType"
	ret, err := internal.OnceCall(key, func() (*[]int, error) {
		ret, err := h.findHandlerFieldIndexFormType(typ)
		return &ret, err
	})
	return *ret, err

}
func (h *helperType) findHandlerFieldIndexFormType(typ reflect.Type) ([]int, error) {
	return h.findHandlerFieldIndexFormTypeInternal(typ, make(map[reflect.Type]struct{}))
}

func (h *helperType) findHandlerFieldIndexFormTypeInternal(typ reflect.Type, visited map[reflect.Type]struct{}) ([]int, error) {

	if services.ServiceUtils.IsInjector(typ) {
		return nil, nil
	}

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if typ == reflect.TypeOf(Handler{}) || typ.ConvertibleTo(reflect.TypeOf(Handler{})) {
		return []int{}, nil
	}
	if typ.Kind() != reflect.Struct {
		return nil, nil
	}

	// Tránh vòng lặp vô hạn
	if _, ok := visited[typ]; ok {
		return nil, nil
	}
	visited[typ] = struct{}{}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldType := field.Type
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}
		if fieldType.Kind() != reflect.Struct {
			continue
		}

		if fieldType == reflect.TypeOf(Handler{}) {
			return []int{i}, nil
		}

		fieldIndex, err := h.findHandlerFieldIndexFormTypeInternal(fieldType, visited)
		if err != nil {
			return nil, err
		}
		if fieldIndex != nil {
			return append([]int{i}, fieldIndex...), nil
		}
	}
	return nil, nil
}
