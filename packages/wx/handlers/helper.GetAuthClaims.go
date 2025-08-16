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
	return h.getAuthClaimsInternal(typ, make(map[reflect.Type]struct{}))
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
