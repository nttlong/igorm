package handlers

import "reflect"

var HandlerIsArgHandler func(typ reflect.Type, visited map[reflect.Type]bool) ([]int, []int, []int)

func (h *helperType) HandlerIsArgHandler(typ reflect.Type, visited map[reflect.Type]bool) ([]int, []int, []int) {
	return HandlerIsArgHandler(typ, visited)
}

func (h *helperType) HandlerFindInMethod(method reflect.Method) ([]int, []int, []int, error) {
	for i := 1; i < method.Type.NumIn(); i++ {
		typ := method.Type.In(i)
		if ret, reqIndex, resIndex := h.HandlerIsArgHandler(typ, make(map[reflect.Type]bool)); ret != nil {
			return ret, reqIndex, resIndex, nil
		}
	}
	return nil, nil, nil, nil

}
