package vapi

import (
	"reflect"
	"strings"
)

func (s *HtttpServer) AddController(controllers ...interface{}) {

}

func InspectMethod[T any]() ([]handlerInfo, error) {
	ret := []handlerInfo{}
	typ := reflect.TypeFor[*T]()
	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)
		info, err := inspector.helper.GetHandlerInfo(method)
		if err != nil {
			return nil, err
		}
		if info == nil {
			continue
		}
		ret = append(ret, *info)
	}
	return ret, nil

}

func Controller[T any](groupTags string, route string, init func() (*T, error)) error {
	list, err := InspectMethod[T]()
	if err != nil {
		return err
	}
	for _, x := range list {

		wHandler := webHandler{
			apiInfo:   x,
			initFunc:  reflect.ValueOf(init),
			routePath: route + "/" + x.UriHandler,
			method:    x.HttpMethod,
		}
		wHandler.routePath = strings.ReplaceAll(wHandler.routePath, "//", "/")
		handlerList = append(handlerList, wHandler)
	}
	return nil
}
