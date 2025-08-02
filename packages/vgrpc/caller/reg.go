package caller

import (
	"encoding/json"
	"fmt"
	reflect "reflect"
	"strings"
	sync "sync"
)

type initAddSingletonServiceInfo struct {
	val      typMap
	typePath string
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
			// serviceInstaneType: ServiceInstaneTypes_Singleton,
			funcs:    map[string]funcInfo{},
			instance: reflect.ValueOf(instance),
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

type initAddSingletonService struct {
	once sync.Once
	val  initAddSingletonServiceInfo
	err  error
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

type funcInfo struct {
	input  []reflect.Type
	output []reflect.Type
	method reflect.Method
}
type singletonCallerDataInfo struct {
	fn  funcInfo
	ins reflect.Value
}
type typMap struct {
	// serviceInstaneType ServiceInstaneTypes
	funcs    map[string]funcInfo
	instance reflect.Value
}

var (
	regInfoSingletonData map[string]singletonCallerDataInfo

	// regInfo                 map[string]typMap //<-- package name/struct type name/function name
	addSingletonServiceOnce sync.Map
)

func callNoAgs(methoPath string) ([]byte, error) {
	typemap, ok := regInfoSingletonData[strings.ToLower(methoPath)]
	if !ok {
		return nil, fmt.Errorf("%s was not found", methoPath)
	}

	ret := typemap.fn.method.Func.Call([]reflect.Value{typemap.ins})
	retIns := make([]interface{}, len(ret))
	for i, item := range ret {
		retIns[i] = item.Interface()
	}
	retJson, err := json.Marshal(retIns)
	if err != nil {
		return nil, err
	}

	return retJson, nil

}

func Call(methoPath string, input []byte) ([]byte, error) {
	defer func() {
		// recover() chỉ có tác dụng trong một hàm defer
		if r := recover(); r != nil {
			// In ra thông điệp panic
			fmt.Println("Đã phục hồi từ một panic:", r)
			fmt.Println("loi")
		}
	}()
	if len(input) == 0 {
		return callNoAgs(methoPath)
	}
	typemap, ok := regInfoSingletonData[strings.ToLower(methoPath)]
	if !ok {
		return nil, fmt.Errorf("%s was not found", methoPath)
	}
	inputVals := make([]reflect.Value, typemap.fn.method.Type.NumIn())
	inputVals[0] = typemap.ins
	inputType := typemap.fn.method.Type.In(1)
	val := reflect.New(inputType)

	inputVal := val.Elem().Interface()
	err := json.Unmarshal(input, &inputVal)
	if err != nil {
		return nil, err
	}

	inputVals[1] = val.Elem()

	ret := typemap.fn.method.Func.Call(inputVals)
	retIns := make([]interface{}, len(ret))
	for i, item := range ret {
		retIns[i] = item.Interface()
	}
	retJson, err := json.Marshal(retIns)
	if err != nil {
		return nil, err
	}

	return retJson, nil

}
