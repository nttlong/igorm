package handlers

import (
	"net/http"
	"reflect"
)

var DependIsHttpService func(typ reflect.Type) bool
var CreateServiceContext func(req *http.Request, res http.ResponseWriter) reflect.Value
var DependIsHttpServiceMethodHasContext func(typ reflect.Type) error

func (h *helperType) DependIsHttpServiceMethodHasContext(typ reflect.Type) error {
	return DependIsHttpServiceMethodHasContext(typ)
}
func (h *helperType) DependIsHttpService(typ reflect.Type) bool {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return false
	}
	if DependIsHttpService(typ) {
		return true
	}
	return false
}
