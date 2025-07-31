package vgrpc

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type ServiceInstaneTypes int

const (
	// Bắt đầu từ 0 (giá trị mặc định)
	ServiceInstaneTypes_None ServiceInstaneTypes = iota
	ServiceInstaneTypes_Singleton
	ServiceInstaneTypes_Transision
)

type funcInfo struct {
	input  []reflect.Type
	output []reflect.Type
	method reflect.Method
}
type typMap struct {
	serviceInstaneType ServiceInstaneTypes
	funcs              map[string]funcInfo
	instance           reflect.Value
}

func analyzeFunction(method reflect.Method) funcInfo {
	ret := funcInfo{

		method: method,
		input:  []reflect.Type{},
		output: []reflect.Type{},
	}

	for i := 1; i < method.Type.NumIn(); i++ {
		ret.input = append(ret.input, method.Type.In(i))
	}

	for i := 0; i < method.Type.NumOut(); i++ {
		ret.output = append(ret.output, method.Type.Out(i))
	}

	return ret
}

type singletonCallerDataInfo struct {
	fn  funcInfo
	ins reflect.Value
}

var (
	regInfoSingletonData map[string]singletonCallerDataInfo

	regInfo                 map[string]typMap //<-- package name/struct type name/function name
	addSingletonServiceOnce sync.Map
)

type initAddSingletonServiceInfo struct {
	val      typMap
	typePath string
}
type initAddSingletonService struct {
	once sync.Once
	val  initAddSingletonServiceInfo
	err  error
}

func addSingletonService[T any](init func() (T, error)) (*initAddSingletonServiceInfo, error) {
	typ := reflect.TypeFor[T]()
	_typ := typ
	if _typ.Kind() == reflect.Ptr {
		_typ = _typ.Elem()
	}

	typePath := strings.ToLower(_typ.String())
	instance, err := init()
	if err != nil {
		return nil, err
	}

	ret := &initAddSingletonServiceInfo{
		typePath: typePath,
		val: typMap{
			serviceInstaneType: ServiceInstaneTypes_Singleton,
			funcs:              map[string]funcInfo{},
			instance:           reflect.ValueOf(instance),
		},
	}
	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)

		retMethod := analyzeFunction(method)
		ret.val.funcs[strings.ToLower(method.Name)] = retMethod
		// regInfo[typePath].funcs[strings.ToLower(method.Name)] = retMethod

	}

	return ret, nil

}
func AddSingletonService[T any](init func() (*T, error)) error {
	key := reflect.TypeFor[T]()
	if key.Kind() == reflect.Ptr {
		key = key.Elem()
	}
	actual, _ := addSingletonServiceOnce.LoadOrStore(key, &initAddSingletonService{})
	item := actual.(*initAddSingletonService)
	item.once.Do(func() {

		val, err := addSingletonService(init)

		item.val = *val
		item.err = err
		if err == nil {
			if regInfo == nil {
				regInfo = map[string]typMap{}
			}
			regInfo[item.val.typePath] = item.val.val
		}
		if regInfoSingletonData == nil {
			regInfoSingletonData = map[string]singletonCallerDataInfo{}
		}

		for k, v := range item.val.val.funcs {
			regInfoSingletonData[item.val.typePath+"."+k] = singletonCallerDataInfo{
				fn:  v,
				ins: item.val.val.instance,
			}

		}

	})
	return item.err

}
func Call2(methoPath string, input []interface{}) ([]interface{}, error) {
	typemap := regInfoSingletonData[strings.ToLower(methoPath)]
	inputVals := make([]reflect.Value, len(input)+1)
	inputVals[0] = typemap.ins

	for i, item := range input {
		inputVals[i+1] = reflect.ValueOf(item)
	}
	ret := typemap.fn.method.Func.Call(inputVals)
	retIns := make([]interface{}, len(ret))
	for i, item := range ret {
		retIns[i] = item.Interface()
	}

	return retIns, nil

}
func Call(methoPath string, input []interface{}) ([]interface{}, error) {
	items := strings.Split(methoPath, ".")
	typePath := strings.ToLower(strings.Join(items[0:len(items)-1], "."))
	funcName := strings.ToLower(items[len(items)-1])

	if mapType, ok := regInfo[typePath]; ok {
		if fun, ok := mapType.funcs[funcName]; ok {

			inputVals := make([]reflect.Value, len(input)+1)
			inputVals[0] = mapType.instance

			for i, item := range input {
				inputVals[i+1] = reflect.ValueOf(item)
			}

			ret := fun.method.Func.Call(inputVals)
			retIns := make([]interface{}, len(ret))
			for i, item := range ret {
				retIns[i] = item.Interface()
			}
			return retIns, nil

		} else {
			return nil, fmt.Errorf("%s not found", funcName)

		}

	} else {
		return nil, fmt.Errorf("%s not found", methoPath)
	}

}
