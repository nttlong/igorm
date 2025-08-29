package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
)

type mockType struct {
}
type mockHandler struct {
	Info    *HandlerInfo
	Url     string
	Handler http.HandlerFunc
}

func (mock *mockType) findMethod(handlerType reflect.Type, funcName string) *reflect.Method {
	typ := handlerType
	if typ.Kind() == reflect.Struct {
		typ = reflect.PointerTo(typ)

	}
	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)
		if method.Name == funcName {
			return &method
		}
	}
	return nil
}
func (mock *mockType) AddHandler(handlerType reflect.Type, funcName string) (*mockHandler, error) {
	mt := mock.findMethod(handlerType, funcName)
	if mt == nil {
		return nil, fmt.Errorf("method %s not found in %s", funcName, handlerType.String())
	}
	info, err := Helper.getHandlerInfo(*mt)
	if err != nil {
		return nil, err
	}

	ret := &mockHandler{
		Info: info,
		Url:  info.UriHandler,
	}
	ret.Handler = func(w http.ResponseWriter, r *http.Request) {

		Helper.ReqExec.Invoke(*ret.Info, r, w)
	}
	return ret, nil

}
func (mock *mockType) NewFormRequest() *mockRequest {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	return &mockRequest{
		Body:   body,
		Writer: writer,
		req:    nil,
		files:  make([]*os.File, 0),
	}

}
func (mock *mockType) NewJSONRequest(method, url string, payload any) *mockRequest {

	var buf bytes.Buffer
	if payload != nil {
		if err := json.NewEncoder(&buf).Encode(payload); err != nil {
			return &mockRequest{
				Err: err,
			}
		}
	}
	req, err := http.NewRequest(method, url, &buf)
	if err != nil {
		return &mockRequest{
			Err: err,
		}
	}
	req.Header.Set("Content-Type", "application/json")
	return &mockRequest{

		req: req,
	}

}
