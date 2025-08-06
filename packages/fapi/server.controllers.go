package fapi

import (
	"reflect"
	"strings"
)

func (s *HtttpServer) AddController(controllers ...interface{}) {

}

var cacheController map[reflect.Type]apiMethodInfo

type handlerInfo struct {
	SwaggerRoute string
	routePath    string
	apiInfo      apiMethodInfo
	initFunc     reflect.Value
	groupTags    string
}

var handlerList []handlerInfo

func Controller[T any](groupTags string, route string, init func() (*T, error)) {
	if handlerList == nil {
		handlerList = make([]handlerInfo, 0)
	}
	t := reflect.TypeFor[*T]()
	t2 := reflect.TypeFor[T]()
	listOfApiMethodInfo := api.ExtractAllApiMethodInfo(t)
	for _, apiInfo := range listOfApiMethodInfo {
		routePath := apiInfo.routeHandler
		swaggerRoute := apiInfo.tags.Url
		if route != "" {
			pkg := strings.ToLower(t2.String())
			items := strings.Split(pkg, ".")
			routePath = route + "/" + strings.Join(items, "/") + "/" + routePath
			swaggerRoute = route + "/" + strings.Join(items, "/") + "/" + apiInfo.tags.Url
		} else if apiInfo.IsAbsUri() {
			// routePath = apiInfo.tags.Url
			swaggerRoute = apiInfo.tags.Url
		} else {
			pkg := strings.ToLower(t2.String())
			items := strings.Split(pkg, ".")
			routePath = "/" + strings.Join(items, "/") + "/" + routePath
			swaggerRoute = "/" + strings.Join(items, "/") + "/" + apiInfo.tags.Url
		}
		handler := handlerInfo{
			SwaggerRoute: swaggerRoute,
			groupTags:    groupTags,
			routePath:    routePath,
			apiInfo:      apiInfo,
			initFunc:     reflect.ValueOf(init),
		}
		handlerList = append(handlerList, handler)
	}

}
