package handlers

import (
	"net/http"
	"reflect"
	"wx/internal"
)

func (h *helperType) skipTypeWhenGetAuthClaims(fieldType reflect.Type) bool {
	key := fieldType.String() + "/helperType/skipTyoeWhenGetAuthClaims"
	ret, _ := internal.OnceCall(key, func() (*bool, error) {
		ret := false
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}
		if fieldType.Kind() == reflect.Interface {
			ret = true
		}
		if fieldType == reflect.TypeOf(http.Request{}) {
			ret = true
		}
		return &ret, nil
	})
	return *ret
}
func (h *helperType) GetAuthClaims(typ reflect.Type) []int {
	key := typ.String() + "/helperType/GetAuthClaims"
	r, _ := internal.OnceCall(key, func() (*[]int, error) {
		ret := h.getAuthClaimsInternal(typ, make(map[reflect.Type]struct{}))

		return &ret, nil
	})
	return *r
}
func (h *helperType) GetUserClaims(typ reflect.Type) [][]int {
	key := typ.String() + "/helperType/GetUserClaims"
	r, _ := internal.OnceCall(key, func() (*[][]int, error) {
		ret := h.getUserClaimsWithVisited(typ, make(map[reflect.Type]string))
		return &ret, nil
	})
	return *r
}
func (h *helperType) getUserClaimsWithVisited(typ reflect.Type, visited map[reflect.Type]string) [][]int {
	ret := [][]int{}
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil
	}
	if _, ok := visited[typ]; ok {

		return nil
	}
	visited[typ] = typ.String()
	if typ.ConvertibleTo(reflect.TypeOf(UserClaims{})) {
		return ret
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldType := field.Type
		indexes := h.getUserClaimsWithVisited(fieldType, visited)

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

func (h *helperType) getAuthClaimsInternal(typ reflect.Type, visited map[reflect.Type]struct{}) []int {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil
	}

	// Tránh vòng lặp vô hạn
	if _, ok := visited[typ]; ok {
		return nil
	}
	visited[typ] = struct{}{}

	if typ.ConvertibleTo(reflect.TypeOf(UserClaims{})) {
		return []int{}
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldType := field.Type

		if h.skipTypeWhenGetAuthClaims(fieldType) {
			continue
		}

		ret := h.getAuthClaimsInternal(fieldType, visited)
		if ret != nil {

			return append(field.Index, ret...)
		}
	}
	return nil
}
