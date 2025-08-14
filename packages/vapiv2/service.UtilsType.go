package vapi

import (
	"reflect"
	"strings"
)

type serviceUtilsType struct {
	pkgPath                string
	checkSingletonTypeName string
	checkTransientTypeName string
}

var serviceUtils = &serviceUtilsType{
	pkgPath:                reflect.TypeOf(serviceUtilsType{}).PkgPath(),
	checkSingletonTypeName: strings.Split(reflect.TypeOf(Singleton[any, any]{}).String(), "[")[0] + "[",
	checkTransientTypeName: strings.Split(reflect.TypeOf(Transient[any, any]{}).String(), "[")[0] + "[",
}

func (svc *serviceUtilsType) IsFieldSingleton(field reflect.StructField) bool {
	if field.Type.PkgPath() != svc.pkgPath {
		return false
	}
	return strings.HasPrefix(field.Type.String(), svc.checkSingletonTypeName)
}
func (svc *serviceUtilsType) IsFieldTransient(field reflect.StructField) bool {
	if field.Type.PkgPath() != svc.pkgPath {
		return false
	}
	return strings.HasPrefix(field.Type.String(), svc.checkTransientTypeName)

}
func (svc *serviceUtilsType) IsSingletonType(typ reflect.Type) bool {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.PkgPath() != svc.pkgPath {
		return false
	}

	return strings.HasPrefix(typ.String(), svc.checkSingletonTypeName)
}
func (svc *serviceUtilsType) CreateSingeton(receiverValue *reflect.Value, field reflect.StructField) {
	fieldValue := receiverValue.Elem().FieldByIndex(field.Index)
	if fieldValue.Kind() == reflect.Ptr {
		fieldValue = fieldValue.Elem()
	}
	fieldType := field.Type
	if fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
	}

	instanceOfField := reflect.New(fieldType)

	if fieldValue.Kind() == reflect.Ptr {
		fieldValue = fieldValue.Elem()
	}

	fieldValue.Set(instanceOfField.Elem())
}
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
		if svc.IsFieldSingleton(field) || svc.IsFieldTransient(field) {
			return true
		}
		if svc.isInjectorInternal(fieldType, visited) {
			return true
		}
	}
	return false
}

// func (svc *serviceUtilsType) IsInjector(typ reflect.Type) bool {

// 	if typ.Kind() == reflect.Ptr {
// 		typ = typ.Elem()
// 	}
// 	if typ.Kind() != reflect.Struct {
// 		return false
// 	}
// 	for i := 0; i < typ.NumField(); i++ {
// 		field := typ.Field(i)
// 		fieldType := field.Type
// 		if fieldType.Kind() == reflect.Ptr {
// 			fieldType = fieldType.Elem()
// 		}
// 		if fieldType.Kind() != reflect.Struct {
// 			continue
// 		}
// 		if svc.IsFieldSingleton(field) || svc.IsFieldTransient(field) {
// 			return true
// 		}
// 		if svc.IsInjector(fieldType) {
// 			return true
// 		}

// 	}
// 	return false
// }
