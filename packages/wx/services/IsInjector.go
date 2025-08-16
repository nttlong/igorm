/*
This file declare function check a type is inject to follow rule
*/
package services

import "reflect"

func (svc *serviceUtilsType) IsInjector(typ reflect.Type) bool {
	return svc.isInjectorInternal(typ, make(map[reflect.Type]struct{}))
}

func (svc *serviceUtilsType) isInjectorInternal(typ reflect.Type, visited map[reflect.Type]struct{}) bool {
	
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return false
	}

	// Nếu đã kiểm tra rồi thì bỏ qua để tránh vòng lặp
	if _, ok := visited[typ]; ok {
		return false
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
		if svc.IsFieldSingleton(field) || svc.IsFieldScoped(field) {
			return true
		}
		if svc.isInjectorInternal(fieldType, visited) {
			return true
		}
	}
	return false
}
