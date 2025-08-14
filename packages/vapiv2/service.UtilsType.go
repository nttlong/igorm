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
