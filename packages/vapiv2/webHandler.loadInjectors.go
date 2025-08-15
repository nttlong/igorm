package vapi

import (
	"net/http"
	"reflect"
)

func (web *webHandlerRunnerType) LoadInjector(handler webHandler, req *http.Request, res http.ResponseWriter) ([]reflect.Value, error) {
	if len(handler.apiInfo.IndexOfInjectors) == 0 {
		return nil, nil
	}
	ret := make([]reflect.Value, len(handler.apiInfo.IndexOfInjectors))
	for i, injector := range handler.apiInfo.IndexOfInjectors {
		injectorType := handler.apiInfo.Method.Type.In(injector)

		valOfInjector, err := serviceUtils.NewService(injectorType, req, res)
		if err != nil {
			return nil, err
		}
		ret[i] = *valOfInjector

	}
	return ret, nil

}
