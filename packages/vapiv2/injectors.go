package vapi

import (
	"net/http"
	"reflect"
)

type Inject[T any] struct {
	ins     T
	Err     error
	Context HttpContext
}
type HttpContext struct {
	Req *http.Request
	Res http.ResponseWriter
}

var cacheInjectorResolve map[reflect.Type]reflect.Value

func InjectorResolve[T any](fn func(injector ...any) (*T, error)) {
	cacheInjectorResolve[reflect.TypeFor[T]()] = reflect.ValueOf(fn)

}
