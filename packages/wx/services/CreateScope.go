package services

import "reflect"

func (svc *serviceUtilsType) CreateScope(typ reflect.Type) reflect.Value {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return reflect.New(typ)
}
