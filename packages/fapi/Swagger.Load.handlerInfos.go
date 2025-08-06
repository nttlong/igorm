package fapi

import (
	"fapi/swaggers"
	"fmt"
	"reflect"
	"strings"
)

func LoadHandlerInfo(s *swaggers.Swagger) {
	for _, handler := range handlerList {
		op := handler.apiInfo.CreatePathItemFromhandlerInfo()
		op.Tags = append(op.Tags, handler.groupTags)
		pathItem := swaggers.PathItem{}
		pathItemVal := reflect.ValueOf(&pathItem).Elem()

		fieldInfo := pathItemVal.FieldByNameFunc(func(s string) bool {
			return strings.EqualFold(s, handler.apiInfo.httpMethod)

		})
		fieldInfo.Set(reflect.ValueOf(op))

		s.Paths[handler.SwaggerRoute] = pathItem
		fmt.Println(handler.SwaggerRoute)

	}
}
