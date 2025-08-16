package htttpserver

// import (
// 	"reflect"
// 	"strings"
// 	swaggers3 "wx/swagger3"
// )

// func (sb *SwaggerBuild) swagger3GetPaths() *SwaggerBuild {
// 	ret := map[string]swaggers3.PathItem{}

// 	for _, h := range handlerList {
// 		swaggerUri := strings.TrimPrefix(h.apiInfo.Uri, "/")

// 		pathItem := swaggers3.PathItem{}
// 		pathItemType := reflect.TypeOf(pathItem)

// 		fieldHttpMethod, ok := pathItemType.FieldByNameFunc(func(s string) bool {
// 			return strings.EqualFold(s, h.method)
// 		})
// 		if !ok {
// 			continue
// 		}

// 		operation := sb.createOperation(h)
// 		operationValue := reflect.ValueOf(operation)

// 		pathItemValue := reflect.ValueOf(&pathItem).Elem() // lấy địa chỉ struct để set

// 		fieldValue := pathItemValue.FieldByIndex(fieldHttpMethod.Index)
// 		if fieldValue.Kind() == reflect.Ptr {
// 			fieldValue.Set(operationValue) // <<--panic: reflect.Set: value of type swaggers3.Operation is not assignable to type *swaggers3.Operation

// 		} else {
// 			fieldValue.Set(operationValue.Elem())
// 		}

// 		ret["/"+swaggerUri] = pathItem
// 	}

// 	sb.swagger.Paths = ret
// 	return sb
// }
