package vapi

import (
	"fmt"
	"net/http"
	"reflect"
)

func (web *webHandlerRunnerType) ResolveReceiverValue(handler webHandler) (reflect.Value, error) {
	key := handler.apiInfo.ReceiverType.String() + "/webHandlerRunnerType/ResolveReceiverValue"
	ret, err := OnceCall(key, func() (*reflect.Value, error) {
		result := handler.initFunc.Call([]reflect.Value{})
		if result[1].IsValid() && !result[1].IsNil() {
			return nil, result[1].Interface().(error)
		}

		return &result[0], nil
	})
	return *ret, err
}
func (web *webHandlerRunnerType) CreateHttpContext(handler webHandler, w http.ResponseWriter, r *http.Request) (reflect.Value, error) {
	context := reflect.New(handler.apiInfo.TypeOfArgsElem)

	context.Elem().FieldByName("Req").Set(reflect.ValueOf(r))
	context.Elem().FieldByName("Res").Set(reflect.ValueOf(w))
	if handler.apiInfo.IsRegexHandler {
		placeHolders := handler.apiInfo.RegexUriFind.FindStringSubmatch(r.URL.Path)
		fmt.Println(handler.apiInfo.RegexUriFind.String())
		fmt.Println(placeHolders)
		if len(placeHolders) == 0 {
			return context, fmt.Errorf("invalid uri")
		}
		for i, uriParam := range handler.apiInfo.UriParams {
			fieldIndex := uriParam.FieldIndex
			context.Elem().FieldByIndex(fieldIndex).Set(reflect.ValueOf(placeHolders[i+1]))
		}

	}
	return context, nil
}
