package vapi

import "reflect"

func (web *webHandlerRunnerType) MethodCall(handler webHandler, args []reflect.Value) []reflect.Value {
	for i := 0; i < handler.apiInfo.Method.Type.NumIn(); i++ {
		argTyp := handler.apiInfo.Method.Type.In(i)
		if argTyp.Kind() != reflect.Ptr && args[i].Kind() == reflect.Ptr {
			args[i] = args[i].Elem()
		}

	}

	retArgs := handler.apiInfo.Method.Func.Call(args)
	return retArgs
}
