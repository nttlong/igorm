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
	args[handlerInfo.IndexOfArg] = reflect.New(handlerInfo.TypeOfArgsElem)
	if handlerInfo.IndexOfRequestBody != -1 {
		bodyValue, err := reqExec.GetBodyValue(handlerInfo, r)
		if err != nil {
			return nil, err
		}
		args[handlerInfo.IndexOfRequestBody] = *bodyValue

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
