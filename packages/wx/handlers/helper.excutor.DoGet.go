package handlers

import (
	"net/http"
	"reflect"
	wxErrors "wx/errors"
)

func (reqExec *RequestExecutor) LoadInjectors(handlerInfo HandlerInfo, r *http.Request, w http.ResponseWriter) ([]reflect.Value, error) {
	injectors := make([]reflect.Value, len(handlerInfo.IndexOfInjectors))
	for i, x := range handlerInfo.IndexOfInjectors {
		r, err := Helper.DependNew(handlerInfo.Method.Type.In(x))
		if err != nil {
			return nil, err
		}
		injectors[i] = *r
	}
	return injectors, nil
}
func (reqExec *RequestExecutor) LoadInjectorsToArgs(handlerInfo HandlerInfo, r *http.Request, w http.ResponseWriter, args []reflect.Value) error {
	if len(handlerInfo.IndexOfInjectors) > 0 {
		injectors, err := reqExec.LoadInjectors(handlerInfo, r, w)
		if err != nil {
			return err
		}
		for i, x := range handlerInfo.IndexOfInjectors {
			if args[x].Kind() == reflect.Ptr {
				args[x] = injectors[i]
			} else {
				args[x] = injectors[i]
			}

		}

	}
	return nil
}
func (reqExec *RequestExecutor) DoGet(handlerInfo HandlerInfo, r *http.Request, w http.ResponseWriter) (any, error) {
	ctlValue, err := reqExec.CreateControllerValue(handlerInfo)
	if err != nil {
		return nil, wxErrors.NewServiceInitError(err.Error())
	}
	controllerValue := *ctlValue

	args := make([]reflect.Value, handlerInfo.Method.Func.Type().NumIn())
	args[0] = controllerValue
	ctxHandler, err := reqExec.CreateHandlerContext(handlerInfo, r, w)
	if err != nil {
		return nil, err
	}
	args[handlerInfo.IndexOfArg] = *ctxHandler
	err = reqExec.LoadInjectorsToArgs(handlerInfo, r, w, args)
	if err != nil {
		return nil, err
	}
	err = reqExec.LoadInjectorInjectServiceToArgs(handlerInfo, r, w, args)
	if err != nil {
		return nil, err
	}
	if handlerInfo.IndexOfAuthClaimsArg != -1 {
		AuthClaimsType := handlerInfo.Method.Type.In(handlerInfo.IndexOfAuthClaimsArg)
		AuthClaimsValue, err := Helper.DepenAuthCreate(AuthClaimsType, r, w)
		if err != nil {
			return nil, err
		}
		if AuthClaimsType.Kind() == reflect.Ptr {
			args[handlerInfo.IndexOfAuthClaimsArg] = *AuthClaimsValue
		} else {
			args[handlerInfo.IndexOfAuthClaimsArg] = (*AuthClaimsValue).Elem()
		}

	}

	//reqExec.CreateHandler(handlerInfo)
	rets := handlerInfo.Method.Func.Call(args)
	if len(rets) == 0 {
		return nil, nil
	}
	if len(rets) > 1 {
		if err, ok := rets[len(rets)-1].Interface().(error); ok {
			return nil, err
		}
	}
	return rets[0].Interface(), nil

}
func (reqExec *RequestExecutor) LoadInjectorInjectServiceToArgs(handlerInfo HandlerInfo, r *http.Request, w http.ResponseWriter, args []reflect.Value) error {
	if len(handlerInfo.IndexOfInjectorService) > 0 {
		httpContextService := CreateServiceContext(r, w)
		for _, x := range handlerInfo.IndexOfInjectorService {
			injectServiceType := handlerInfo.Method.Type.In(x)

			if injectServiceType.Kind() == reflect.Ptr {
				injectServiceType = injectServiceType.Elem()
			}
			injectServiceValue := reflect.New(injectServiceType)
			httpContextField := injectServiceValue.Elem().FieldByName("HttpContext")
			if httpContextField.IsValid() && httpContextField.CanSet() {
				httpContextField.Set(httpContextService)
			}

			if handlerInfo.Method.Type.In(x).Kind() == reflect.Ptr {
				args[x] = injectServiceValue
			} else {
				args[x] = injectServiceValue.Elem()
			}

		}
	}

	return nil
}
