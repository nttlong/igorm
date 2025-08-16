package handlers

import (
	"fmt"
	"mime/multipart"
	"reflect"
)

func (h *helperType) IsFieldFileUpload(field reflect.StructField) bool {
	fmt.Println(field.Name)
	fieldType := field.Type
	if fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
	}
	if fieldType.Kind() == reflect.Slice {
		fieldType = fieldType.Elem()
	}
	if fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
	}
	if reflect.TypeOf(multipart.FileHeader{}) == fieldType {
		return true
	}
	checkType := reflect.TypeOf((*multipart.File)(nil)).Elem()
	if field.Type.Kind() == reflect.Interface && field.Type.Implements(checkType) {
		return true
	}
	return false

}
func (h *helperType) FindFormUploadInType(typ reflect.Type) []int {
	return h.findFormUploadInTypeInternal(typ, nil, make(map[reflect.Type]struct{}))
}

func (h *helperType) findFormUploadInTypeInternal(typ reflect.Type, parentIndex []int, visited map[reflect.Type]struct{}) []int {
	ret := []int{}

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return ret
	}

	// Tránh vòng lặp vô hạn
	if _, ok := visited[typ]; ok {
		return ret
	}
	visited[typ] = struct{}{}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldType := field.Type
		indexPath := append(parentIndex, field.Index...)

		if h.IsFieldFileUpload(field) {
			ret = append(ret, indexPath...)
			continue
		}

		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}
		if fieldType.Kind() == reflect.Struct {
			ret = append(ret, h.findFormUploadInTypeInternal(fieldType, indexPath, visited)...)
		}
	}

	return ret
}
