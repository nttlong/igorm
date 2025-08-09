package vapi

import (
	"reflect"
	"strings"
)

func (s *HtttpServer) AddController(controllers ...interface{}) {

}

// var cacheController map[reflect.Type]handlerInfo

// type handlerInfoSwagger struct {
// 	SwaggerRoute string
// 	routePath    string
// 	apiInfo      handlerInfo
// 	initFunc     reflect.Value
// 	groupTags    string
// }

// var handlerInfoSwaggers []handlerInfoSwagger = []handlerInfoSwagger{}
type webHandler struct {
	routePath string
	apiInfo   handlerInfo
	initFunc  reflect.Value
	method    string
}

var handlerList []webHandler = []webHandler{}

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
		if x.UriHandler[len(x.UriHandler)-1] == '$' {
			panic("loi")
		}
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
