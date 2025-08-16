package handlers

import (
	"fmt"
	"reflect"
	"wx/internal"
)

func (h *helperType) ExtractTags(typ reflect.Type, fieldIndex []int) []string {

	if len(fieldIndex) == 0 {
		return nil
	}

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	key := typ.String() + "/helperType/ExtractTags"
	for _, x := range fieldIndex {
		key += "/" + fmt.Sprint(x)
	}
	ret, _ := internal.OnceCall(key, func() (*[]string, error) {
		ret := []string{}
		field := typ.FieldByIndex([]int{fieldIndex[0]})

		ret = append(ret, field.Tag.Get("route"))
		fieldType := field.Type

		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}
		subRet := h.ExtractTags(fieldType, fieldIndex[1:])

		ret = append(ret, subRet...)

		return &ret, nil
	})
	return *ret

}
