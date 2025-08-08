package vapi

import (
	"fmt"
	"reflect"
	"strings"
)

type handlerInfo struct {
	IndexOfArg int
	FieldIndex []int
}

func handlerInfoFindHandlerFieldIndexFormType(typ reflect.Type) ([]int, error) {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

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
		fieldIndex, err := handlerInfoFindHandlerFieldIndexFormType(fieldType)
		if err != nil {
			return nil, err
		}
		if fieldIndex != nil {
			return append([]int{i}, fieldIndex...), nil
		}
	}
	return nil, nil
}
func handlerInfoFromMethod(method reflect.Method) (*handlerInfo, error) {
	for i := 0; i < method.Type.NumIn(); i++ {
		typ := method.Type.In(i)
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		if typ.Kind() != reflect.Struct {
			continue
		}
		if typ == reflect.TypeOf(Handler{}) {
			return &handlerInfo{
				IndexOfArg: i,
				FieldIndex: []int{},
			}, nil
		}
		fieldIndex, err := handlerInfoFindHandlerFieldIndexFormType(typ)
		if err != nil {
			return nil, err
		}
		if fieldIndex != nil {
			return &handlerInfo{
				IndexOfArg: i,
				FieldIndex: fieldIndex,
			}, nil
		}
	}
	return nil, nil
}
func handlerInfoExtractTags(typ reflect.Type, fieldIndex []int) []string {
	if len(fieldIndex) == 0 {
		return nil
	}

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	ret := []string{}
	field := typ.FieldByIndex([]int{fieldIndex[0]})
	fmt.Println(field.Name)
	ret = append(ret, field.Tag.Get("route"))
	fieldType := field.Type

	if fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
	}
	subRet := handlerInfoExtractTags(fieldType, fieldIndex[1:])

	ret = append(ret, subRet...)

	return ret

}
func handlerInfoExtractUriFromTags(tags []string) string {
	ret := ""
	for i := len(tags) - 1; i >= 0; i-- {
		tag := tags[i]
		if tag == "" {
			continue
		}
		items := strings.Split(tag, ";")
		for _, item := range items {
			if strings.HasPrefix(item, "uri:") {

				val := item[4:]

				if val != "" {
					if strings.Contains(ret, "@") {
						ret = strings.Replace(ret, "@", val, 1)
					} else {
						ret += "/" + val
					}
				}

			}
		}
	}

	return strings.TrimPrefix(strings.TrimSuffix(ret, "/"), "/")

}
