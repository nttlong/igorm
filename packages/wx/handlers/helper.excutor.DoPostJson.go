package handlers

import (
	"net/http"
	"reflect"
	wxErrors "wx/errors"
)

func (reqExec *RequestExecutor) DoJsonPost(handlerInfo HandlerInfo, r *http.Request, w http.ResponseWriter) (interface{}, error) {
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
	err = reqExec.LoadInjectorInjectServiceToArgs(handlerInfo, r, w, args)
	if err != nil {
		return nil, err
	}
	err = reqExec.LoadInjectorsToArgs(handlerInfo, r, w, args)
	if err != nil {
		return nil, err
	}
	args[handlerInfo.IndexOfArg] = *ctxHandler
	if handlerInfo.IndexOfRequestBody != -1 {
		bodyValue, err := reqExec.GetBodyValue(handlerInfo, r)
		if err != nil {
			return nil, err
		}
		if handlerInfo.TypeOfRequestBody.Kind() == reflect.Ptr {

			args[handlerInfo.IndexOfRequestBody] = *bodyValue
		} else {
			args[handlerInfo.IndexOfRequestBody] = (*bodyValue).Elem()
		}

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
	if len(rets) > 0 {
		if err, ok := rets[len(rets)-1].Interface().(error); ok {
			return nil, err
		}
	}
	return rets[0].Interface(), nil

}
