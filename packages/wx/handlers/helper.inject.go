package handlers

import (
	"net/http"
	"reflect"
)

type ServiceType struct {
}

var IsServiceContext func(typ reflect.Type) bool

func (inject *ServiceType) IsServiceContext(typ reflect.Type) bool {
	return IsServiceContext(typ)

}

var GetNewMethodOfServiceContext func(typ reflect.Type) (*reflect.Method, error)

func (inject *ServiceType) GetNewMethod(typ reflect.Type) (*reflect.Method, error) {
	return GetNewMethodOfServiceContext(typ)

}

var NewInjectByType func(typ reflect.Type, newMethod reflect.Method, r *http.Request, w http.ResponseWriter) (*reflect.Value, error)

func (inject *ServiceType) NewInjectByType(typ reflect.Type, newMethod reflect.Method, r *http.Request, w http.ResponseWriter) (*reflect.Value, error) {
	return NewInjectByType(typ, newMethod, r, w)

}
func (inject *ServiceType) LoadService(handlerInfo HandlerInfo, r *http.Request, w http.ResponseWriter, args []reflect.Value) error {
	for i, x := range handlerInfo.ServiceContextArgs {
		injectType := handlerInfo.Method.Type.In(x)
		injectValue, err := inject.NewInjectByType(injectType, handlerInfo.ServiceContextNewMethods[i], r, w)
		if err != nil {
			return err
		}
		args[x] = *injectValue

	}
	return nil
}
