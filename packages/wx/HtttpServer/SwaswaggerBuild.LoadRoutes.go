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

		// pathItem := swaggers3.PathItem{}

		// // pathItemType := reflect.TypeOf(pathItem)
		// retPaths[k] = pathItem
		// operation := sb.createOperation(WebHandler{
		// 	RoutePath: k,
		// 	ApiInfo:   v.Info,
		// 	Method:    v.Info.HttpMethod,
		// })

	}
	sb.swagger3GetPaths()
	// sb.swagger.Paths = retPaths
	// for _, h := range HandlerList {

	// 	swaggerUri := strings.TrimPrefix(strings.ReplaceAll(h.ApiInfo.Uri, "*", ""), "/")

	// 	pathItem := swaggers3.PathItem{}
	// 	pathItemType := reflect.TypeOf(pathItem)

	// 	fieldHttpMethod, ok := pathItemType.FieldByNameFunc(func(s string) bool {
	// 		return strings.EqualFold(s, h.Method)
	// 	})
	// 	if !ok {
	// 		continue
	// 	}

	// 	operation := sb.createOperation(h)
	// 	operationValue := reflect.ValueOf(operation)

	// 	pathItemValue := reflect.ValueOf(&pathItem).Elem() // lấy địa chỉ struct để set

	// 	fieldValue := pathItemValue.FieldByIndex(fieldHttpMethod.Index)
	// 	if fieldValue.Kind() == reflect.Ptr {
	// 		fieldValue.Set(operationValue) // <<--panic: reflect.Set: value of type swaggers3.Operation is not assignable to type *swaggers3.Operation

	// 	} else {
	// 		fieldValue.Set(operationValue.Elem())
	// 	}

	// 	ret["/"+swaggerUri] = pathItem
	// }

	return sb
}
