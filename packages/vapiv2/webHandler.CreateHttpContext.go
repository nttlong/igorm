package vapi

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	vError "vapi/errors"
)

func (web *webHandlerRunnerType) CreateHttpContext(handler webHandler, w http.ResponseWriter, r *http.Request) (reflect.Value, error) {
	context := reflect.New(handler.apiInfo.TypeOfArgsElem)

	context.Elem().FieldByName("Req").Set(reflect.ValueOf(r))
	context.Elem().FieldByName("Res").Set(reflect.ValueOf(w))
	context.Elem().FieldByName("BaseUrl").Set(reflect.ValueOf(handler.apiInfo.BaseUrl))
	if handler.apiInfo.IsRegexHandler {
		placeHolders := handler.apiInfo.RegexUriFind.FindStringSubmatch(r.URL.Path)
		// fmt.Println(handler.apiInfo.RegexUriFind.String())
		// fmt.Println(placeHolders)
		if len(placeHolders) == 0 {
			return context, fmt.Errorf("invalid uri")
		}
		for i, uriParam := range handler.apiInfo.UriParams {
			//fieldIndex := uriParam.FieldIndex
			fmt.Println(handler.apiInfo.TypeOfArgsElem.String())
			fieldSet, ok := handler.apiInfo.TypeOfArgsElem.FieldByNameFunc(func(s string) bool {
				return strings.EqualFold(s, uriParam.Name)
			})
			if !ok {

				return context, vError.NewUriParamParseError(uriParam.Name, handler.apiInfo.TypeOfArgsElem)
			}

			//fieldSet := context.Elem().FieldByName(uriParam.Name)
			valueSet := reflect.ValueOf(placeHolders[i+1])
			fielValueSet := context.Elem().FieldByIndex(fieldSet.Index)
			// if fieldSet.t.Kind() == reflect.Ptr {
			// 	fieldSet = fieldSet.Elem()
			// }
			// if valueSet.Kind() == reflect.Ptr {
			// 	valueSet = valueSet.Elem()
			// }
			if fielValueSet.CanConvert(valueSet.Type()) {
				fielValueSet.Set(valueSet)
			} else {

				return context, vError.NewUriParamConvertError(uriParam.Name, valueSet.Type(), fielValueSet.Type())

			}

		}

	}
	return context, nil
}
