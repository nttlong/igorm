package vapi

import "reflect"

func (web *webHandlerRunnerType) LoadInjector(handler webHandler) ([]reflect.Value, error) {
	if len(handler.apiInfo.IndexOfInjectors) == 0 {
		return nil, nil
	}
	ret := make([]reflect.Value, len(handler.apiInfo.IndexOfInjectors))
	for i, injector := range handler.apiInfo.IndexOfInjectors {
		injectorType := handler.apiInfo.Method.Type.In(injector)
		ret[i] = reflect.New(injectorType)

	}
	return ret, nil

}
