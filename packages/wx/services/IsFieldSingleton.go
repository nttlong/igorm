package services

import (
	"reflect"
	"strings"
)

func (svc *serviceUtilsType) IsFieldSingleton(field reflect.StructField) bool {
	typ := field.Type
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	// if typ.PkgPath() != svc.pkgPath {
	// 	return false

	// }

	return strings.HasPrefix(typ.String(), svc.checkSingletonTypeName())
}
func (svc *serviceUtilsType) IsFieldScoped(field reflect.StructField) bool {
	typ := field.Type
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	// if typ.PkgPath() != svc.pkgPath {
	// 	return false

	// }

	return strings.HasPrefix(typ.String(), svc.checkScopeTypeName())

}
func (svc *serviceUtilsType) IsSingletonType(typ reflect.Type) bool {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.PkgPath() != svc.pkgPath {
		return false
	}

	return strings.HasPrefix(typ.String(), svc.checkSingletonTypeName())
}
