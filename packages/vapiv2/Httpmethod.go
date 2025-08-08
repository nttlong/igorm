package vapi

import (
	"reflect"
	"strings"
	"sync"
)

type Context struct {
}
type HttpGet struct {
	Context
}
type HttpPost struct {
	Context
}
type HttpPut struct {
	Context
}
type HttpDelete struct {
	Context
}
type HttpPatch struct {
	Context
}
type httpUtilsType struct {
}
type initInspectHttpMethodFromType struct {
	val  string
	once sync.Once
}

var cacheInspectHttpMethodFromType sync.Map

func (u *httpUtilsType) inspectHttpMethodFromType(typ reflect.Type) string {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	switch typ {
	case reflect.TypeOf(HttpGet{}):
		return "get"
	case reflect.TypeOf(HttpPost{}):
		return "post"
	case reflect.TypeOf(HttpPut{}):
		return "put"
	case reflect.TypeOf(HttpDelete{}):
		return "delete"
	case reflect.TypeOf(HttpPatch{}):
		return "patch"
	default:
		if typ.Kind() == reflect.Struct {
			for i := 0; i < typ.NumField(); i++ {
				if ret := u.inspectHttpMethodFromType(typ.Field(i).Type); ret != "" {
					return ret
				}
			}
		}
		return ""
	}
}

/*
Each type have only one time for detection
*/
func (u *httpUtilsType) InspectHttpMethodFromType(typ reflect.Type) string {
	actual, _ := cacheInspectHttpMethodFromType.LoadOrStore(typ, &initInspectHttpMethodFromType{})
	item := actual.(*initInspectHttpMethodFromType)
	item.once.Do(func() {
		item.val = u.inspectHttpMethodFromType(typ)
	})
	return item.val
}
func (u *httpUtilsType) IsInjector(typ reflect.Type) bool {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return false
	}
	if typ.PkgPath() == reflect.TypeOf(httpUtilsType{}).PkgPath() && strings.Contains(typ.String(), "Inject[") {
		return true

	}
	return false
}
func (u *httpUtilsType) GetRouteTag(typ reflect.Type) string {

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	httpMethd := u.InspectHttpMethodFromType(typ)
	if httpMethd == "" {
		return ""
	}
	for i := 0; i < typ.NumField(); i++ {
		if ret := u.InspectHttpMethodFromType(typ.Field(i).Type); ret != "" {
			return typ.Field(i).Tag.Get("route")
		}
	}
	return ""
}

var httpUtilsTypeInstance = httpUtilsType{}
