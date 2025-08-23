package wx

import (
	"reflect"
	"strings"
	httpServer "wx/HtttpServer"
	"wx/handlers"
)

func InspectMethod[T any]() ([]handlers.HandlerInfo, error) {
	ret := []handlers.HandlerInfo{}
	typ := reflect.TypeFor[*T]()
	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)

		info, err := handlers.Helper.GetHandlerInfo(method)
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

func CreateWebHandler[T any](x handlers.HandlerInfo, init func() (*T, error)) httpServer.WebHandler {
	wHandler := httpServer.WebHandler{
		ApiInfo:   x,
		InitFunc:  reflect.ValueOf(init),
		RoutePath: "/" + x.UriHandler,
		Method:    x.HttpMethod,
	}
	return wHandler
}
func LoadController[T any](init func() (*T, error)) error {
	list, err := InspectMethod[T]()
	if err != nil {
		return err
	}
	for _, x := range list {

		wHandler := CreateWebHandler(x, init)
		wHandler.RoutePath = strings.ReplaceAll(wHandler.RoutePath, "//", "/")
		httpServer.HandlerList = append(httpServer.HandlerList, wHandler)
	}
	return nil
}
func Routes(baseUri string, ins ...any) error {
	return handlers.Helper.Routes.Add(baseUri, ins...)
}
