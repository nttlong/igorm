package vapi

import (
	"reflect"
	vapiErr "vapi/errors"
)

func (web *webHandlerRunnerType) MethodCall(handler webHandler, args []reflect.Value) ([]reflect.Value, error) {
	for i := 1; i < handler.apiInfo.Method.Type.NumIn(); i++ {
		argTyp := handler.apiInfo.Method.Type.In(i)
		if argTyp.Kind() != reflect.Ptr {
			if args[i].Kind() == reflect.Ptr {
				args[i] = args[i].Elem()
			} else {
				if i == handler.apiInfo.IndexOfRequestBody {
					return nil, vapiErr.NewBadRequestError("invalid argument, missing body data")
				} else {
					return nil, vapiErr.NewBadRequestError("invalid request")
				}
			}
		}

	}

	retArgs := handler.apiInfo.Method.Func.Call(args)
	return retArgs, nil
}
