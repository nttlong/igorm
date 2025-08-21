package htttpserver

import (
	"wx/handlers"
)

func (sb *SwaggerBuild) LoadFromRoutes() *SwaggerBuild {
	dataRoute := handlers.Helper.Routes.Data
	// ret := map[string]swaggers3.PathItem{}
	// retPaths := map[string]swaggers3.PathItem{}
	if HandlerList == nil {
		HandlerList = []WebHandler{}
	}
	for k, v := range dataRoute {
		swagerInfo := WebHandler{
			RoutePath: k,
			ApiInfo:   v.Info,
			Method:    v.Info.HttpMethod,
		}
		HandlerList = append(HandlerList, swagerInfo)

	}
	sb.swagger3GetPaths()

	return sb
}
