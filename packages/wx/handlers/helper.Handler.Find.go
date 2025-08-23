package handlers

import "reflect"

func (h *helperType) HandlerIsArgHandler(typ reflect.Type, visited map[reflect.Type]bool) []int {
	if typ == nil || visited[typ] {
		return nil
	}
	visited[typ] = true
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ == reflect.TypeOf(Handler{}) || typ.ConvertibleTo(reflect.TypeOf(Handler{})) {
		return []int{}
	}

	if typ.Kind() != reflect.Struct {
		return nil
	}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if indexes := h.HandlerIsArgHandler(field.Type, visited); indexes != nil {
			fieldIndex := append(field.Index, indexes...)
			return fieldIndex
		}
	}
	return nil
}
func (h *helperType) HandlerFindInMethod(method reflect.Method) ([]int, error) {
	for i := 1; i < method.Type.NumIn(); i++ {
		typ := method.Type.In(i)
		if ret := h.HandlerIsArgHandler(typ, make(map[reflect.Type]bool)); ret != nil {
			return ret, nil
		}
	}
	return nil, nil

}
