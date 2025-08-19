package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"sync"
	wxErrors "wx/errors"
)

// MockResponseWriter implements http.ResponseWriter

// Helper function tạo *http.Request mock
func NewMockRequest(method, urlStr string, body io.Reader, query url.Values, headers map[string]string) *http.Request {
	req, _ := http.NewRequest(method, urlStr, body)

	// Thêm query params
	if query != nil {
		req.URL.RawQuery = query.Encode()
	}

	// Thêm headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return req
}

type RequestExecutor struct {
}
type initCreateControllerValue struct {
	Controller reflect.Value
	once       sync.Once
}

var cacheCreateControllerValue sync.Map

func (reqExec *RequestExecutor) CreateControllerValueOnce(handlerInfo HandlerInfo) reflect.Value {
	actual, _ := cacheCreateControllerValue.LoadOrStore(handlerInfo.ReceiverTypeElem, &initCreateControllerValue{})
	item := actual.(*initCreateControllerValue)
	item.once.Do(func() {
		item.Controller = reflect.New(handlerInfo.ReceiverTypeElem)
	})
	return item.Controller

}
func (reqExec *RequestExecutor) ResovleNewMethod(instanceVale reflect.Value, method reflect.Method) error {
	args := make([]reflect.Value, method.Type.NumIn())
	args[0] = instanceVale
	for i := 1; i < method.Type.NumIn(); i++ {
		if Helper.IsGenericDepen(method.Type.In(i)) {
			insVal, err := Helper.DependNewOnce(method.Type.In(i))
			if err != nil {
				return err
			}
			if method.Type.In(i).Kind() == reflect.Ptr {
				args[i] = *insVal
			} else {
				args[i] = (*insVal).Elem()
			}
		}
	}
	rets := method.Func.Call(args)
	if len(rets) == 1 {
		if err, ok := rets[0].Interface().(error); ok {
			return err
		}
	}
	return nil
}
func (reqExec *RequestExecutor) CreateControllerInitAllDepenFields(instanceValue reflect.Value, handlerInfo HandlerInfo) error {
	for i := 0; i < instanceValue.Elem().NumField(); i++ {
		field := instanceValue.Elem().Field(i)
		if field.IsValid() && field.CanSet() {

			if Helper.IsGenericDepen(field.Type()) {
				insVal, err := Helper.DependNewOnce(field.Type())
				if err != nil {
					return err
				}
				if field.Type().Kind() == reflect.Ptr {
					field.Set(*insVal)
				} else {
					field.Set((*insVal).Elem())
				}

			}
		}
	}
	return nil
}
func (reqExec *RequestExecutor) CreateControllerValue(handlerInfo HandlerInfo) (*reflect.Value, error) {

	ret := reflect.New(handlerInfo.ReceiverTypeElem)
	err := reqExec.CreateControllerInitAllDepenFields(ret, handlerInfo)
	if err != nil {
		return nil, err
	}
	for i := 0; i < handlerInfo.ReceiverType.NumMethod(); i++ {
		if handlerInfo.ReceiverType.Method(i).Name == "New" {
			err := reqExec.ResovleNewMethod(ret, handlerInfo.ReceiverType.Method(i))
			if err != nil {
				return nil, err
			}

		}
	}
	return &ret, nil

}
func (reqExec *RequestExecutor) CreateHandler(handlerInfo HandlerInfo) reflect.Value {

	return reflect.New(handlerInfo.TypeOfArgsElem)
}
func (reqExec *RequestExecutor) GetBodyValue(handlerInfo HandlerInfo, r *http.Request) (*reflect.Value, error) {
	bodyData := reflect.New(handlerInfo.TypeOfRequestBodyElem)
	if r.Body != nil {

		if err := json.NewDecoder(r.Body).Decode(bodyData.Interface()); err != nil {

			return nil, err
		}
	} else if handlerInfo.TypeOfRequestBody.Kind() == reflect.Struct {
		return nil, wxErrors.NewBadRequestError("request body is required")
	}

	return &bodyData, nil

}

