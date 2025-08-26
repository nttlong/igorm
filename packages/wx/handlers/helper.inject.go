package handlers

import (
	"net/http"
	"reflect"
)

type InjectType struct {
}

var IsInjectType func(typ reflect.Type) bool

func (inject *InjectType) IsInjectType(typ reflect.Type) bool {
	return IsInjectType(typ)

}

var IsReadyRegister func(typ reflect.Type) bool

func (inject *InjectType) IsReadyRegister(typ reflect.Type) bool {
	return IsReadyRegister(typ)

}

var NewInjectByType func(typ reflect.Type, r *http.Request, w http.ResponseWriter) (reflect.Value, error)

func (inject *InjectType) NewInjectByType(typ reflect.Type, r *http.Request, w http.ResponseWriter) (reflect.Value, error) {
	return NewInjectByType(typ, r, w)

}
func (inject *InjectType) LoadInject(handlerInfo HandlerInfo, r *http.Request, w http.ResponseWriter, args []reflect.Value) error {
	for _, x := range handlerInfo.IndexOfArgIsInject {
		injectType := handlerInfo.Method.Type.In(x)
		injectValue, err := inject.NewInjectByType(injectType, r, w)
		if err != nil {
			return err
		}
		args[x] = injectValue

	}
	return nil
}
