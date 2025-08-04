package vapi

import (
	"fmt"
	"reflect"
	"strings"
)

var (
	mapHandler                 map[string]reflect.Method
	mapInstanceInit            map[string]reflect.Value
	mapInstanceType            map[string]reflect.Type
	httpMethodMap              map[string]string
	mapIndexOfUserField        map[string][]int
	mapIndexOfHttpContextField map[string][]int
	mapInputParamInfo          map[string]inputParamInfo
)

func AddHandler[T any](initHanlder func() (*T, error)) error {

	swaggerData, err := loadSwaggerInfo()
	if err != nil {
		return err
	}

	typ := reflect.TypeFor[T]()
	typePrt := reflect.TypeFor[*T]()

	apiPath := typ.String()
	if !strings.HasSuffix(apiPath, "Controller") {
		return fmt.Errorf("%s not end with Controller", apiPath)
	}

	indexOfHttpContextField := getHttpContextFieldIndex(typ)

	apiPath = strings.TrimSuffix(apiPath, "Controller")
	apiPath = strings.ToLower(apiPath)
	apiPath = strings.ReplaceAll(apiPath, ".", "/")

	mapUrlPath := swaggerData.BasePath + "/" + apiPath
	for i := 0; i < typePrt.NumMethod(); i++ {
		method := typePrt.Method(i)
		methodName := method.Name
		if !strings.HasSuffix(methodName, "_Get") &&
			!strings.HasSuffix(methodName, "_Post") &&
			!strings.HasSuffix(methodName, "_Put") &&
			!strings.HasSuffix(methodName, "_Delete") &&
			!strings.HasSuffix(methodName, "_Patch") {
			continue

		}
		items := strings.Split(methodName, "_")
		httpMethod := strings.ToUpper(items[len(items)-1])

		methodName = strings.Join(items[0:len(items)-1], "_")
		methodName = strings.ToLower(methodName)

		mapHandler[mapUrlPath+"/"+methodName] = method
		mapInstanceType[mapUrlPath+"/"+methodName] = typePrt

		mapInstanceInit[mapUrlPath+"/"+methodName] = reflect.ValueOf(initHanlder)
		httpMethodMap[mapUrlPath+"/"+methodName] = httpMethod
		mapIndexOfHttpContextField[mapUrlPath+"/"+methodName] = indexOfHttpContextField

		fullApi := apiPath + "/" + methodName
		pathItem, inputParamInfo := getPathItem(httpMethod, method, typ)
		mapInputParamInfo[mapUrlPath+"/"+methodName] = inputParamInfo

		swaggerData.Paths["/"+fullApi] = pathItem
	}
	return nil

}
func init() {
	mapHandler = map[string]reflect.Method{}
	mapInstanceInit = map[string]reflect.Value{}
	mapInstanceType = map[string]reflect.Type{}
	httpMethodMap = map[string]string{}
	mapIndexOfUserField = map[string][]int{}
	mapIndexOfHttpContextField = map[string][]int{}
	mapInputParamInfo = map[string]inputParamInfo{}

}
