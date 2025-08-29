package handlers

import (
	"net/http"
	"reflect"
	"strings"
	wxErrors "wx/errors"
	"wx/internal"
)

func (reqExec *RequestExecutor) GetFieldByName(typ reflect.Type, fieldName string) *reflect.StructField {
	key := typ.String() + "/RequestExecutor/GetFieldByName/" + fieldName
	ret, err := internal.OnceCall(key, func() (*reflect.StructField, error) {
		ret, ok := typ.FieldByNameFunc(func(s string) bool {
			return strings.EqualFold(s, fieldName)
		})
		if !ok {
			return nil, nil
		}
		return &ret, nil
	})
	if err != nil {
		return nil
	}
	return ret
}

// type formBodyItem struct {
// 	IndexFields [][]int
// 	Value       interface{}
// 	IsRequire   bool
// }

func (reqExec *RequestExecutor) DoFormPost(handlerInfo HandlerInfo, r *http.Request, w http.ResponseWriter) (interface{}, error) {
	ctlValue, err := reqExec.CreateControllerValue(handlerInfo)
	if err != nil {
		return nil, wxErrors.NewServiceInitError(err.Error())
	}
	controllerValue := *ctlValue

	args := make([]reflect.Value, handlerInfo.Method.Func.Type().NumIn())
	args[0] = controllerValue
	args[handlerInfo.IndexOfArg] = reflect.New(handlerInfo.TypeOfArgsElem)
	if handlerInfo.IndexOfRequestBody != -1 {
		bodyValue, err := reqExec.GetFormValue(handlerInfo, r)
		if err != nil {
			return nil, err
		}

		if handlerInfo.Method.Type.In(handlerInfo.IndexOfRequestBody).Kind() == reflect.Ptr {
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
	err = reqExec.LoadInjectorInjectServiceToArgs(handlerInfo, r, w, args)
	if err != nil {
		return nil, err
	}
	err = reqExec.LoadInjectorsToArgs(handlerInfo, r, w, args)
	if err != nil {
		return nil, err
	}
	err = Helper.Services.LoadService(handlerInfo, r, w, args)
	if err != nil {
		return nil, err
	}
	//reqExec.CreateHandler(handlerInfo)
	rets := handlerInfo.Method.Func.Call(args)
	if len(rets) == 0 {
		return nil, nil
	}
	if len(rets) > 0 {
		if err, ok := rets[len(rets)-1].Interface().(error); ok {
			return nil, err
		}
	}
	return rets[0].Interface(), nil

}
