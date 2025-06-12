package dynacall

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type Caller struct {
	Path string
}
type CallerEntry struct {
	Caller     *Caller
	Method     reflect.Method
	FullPath   string
	PkgPath    string
	CallerPath string
}
type ArgsCaller struct {
}

func GetInputTypeOfCallerPath(callerPath string) ([]reflect.Type, error) {
	if !strings.Contains(callerPath, "@") {
		return nil, CallError{
			Err:  fmt.Errorf("caller path should be in format of package.path@method"),
			Code: CallErrorCodeInvalidCallerPath,
		}
	}
	callerEntry, found := callerCache.Load(strings.ToLower(callerPath))
	if !found {
		return nil, CallError{
			Err:  fmt.Errorf("caller not found"),
			Code: CallErrorCodeCallerPathNotFound,
		}
	}
	method := callerEntry.(CallerEntry).Method
	ret := make([]reflect.Type, method.Type.NumIn()-1)
	for i := 1; i < method.Type.NumIn(); i++ {
		ret[i-1] = method.Type.In(i)
	}
	return ret, nil
}

func Call(callerPath string, args interface{}, injector interface{}) (interface{}, error) {
	if !strings.Contains(callerPath, "@") {
		return nil, CallError{
			Err:  fmt.Errorf("caller path should be in format of package.path@method"),
			Code: CallErrorCodeInvalidCallerPath,
		}
	}
	callerEntry, found := callerCache.Load(strings.ToLower(callerPath))
	if !found {
		return nil, CallError{
			Err:  fmt.Errorf("caller not found"),
			Code: CallErrorCodeCallerPathNotFound,
		}
	}
	if caller, ok := callerEntry.(CallerEntry); ok {
		method := caller.Method
		typ := reflect.TypeOf(args)
		fmt.Print(typ.Kind())
		return invoke(method, args, injector)
	}
	return nil, CallError{
		Err:  fmt.Errorf("caller not found"),
		Code: CallerSystemError,
	}

}

var callerCache sync.Map

func getCaller(caller interface{}) *Caller {
	typ := reflect.TypeOf(caller)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if field, found := typ.FieldByName("Caller"); found {
		if field.Type == reflect.TypeOf(Caller{}) {
			val := reflect.ValueOf(caller)
			if val.Kind() == reflect.Ptr {
				val = val.Elem()
			}
			iVal := val.FieldByName(field.Name).Interface()
			ret := iVal.(Caller)
			return &ret
		}

	}
	return nil

}
func getAllMethods(caller interface{}) []reflect.Method {
	ret := []reflect.Method{}
	typ := reflect.TypeOf(caller)
	// if typ.Kind() == reflect.Ptr {
	// 	typ = typ.Elem()
	// }
	fmt.Print(typ.Name())
	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)
		//check if first letter is uppercase
		if method.Name[0] >= 'A' && method.Name[0] <= 'Z' {
			ret = append(ret, method)
		}
	}

	return ret
}
func RegisterCaller(caller interface{}) {
	typ := reflect.TypeOf(caller)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	pkgPath := typ.PkgPath()
	pkgPath = strings.ReplaceAll(pkgPath, "/", ".")
	if pkgPath == "" {
		pkgPath = "main"
	}
	info := getCaller(caller)
	if info != nil {
		methodList := getAllMethods(caller)
		for _, method := range methodList {
			key := method.Name + "@" + pkgPath
			key = strings.ToLower(key)
			callerEntry := CallerEntry{
				Caller:     info,
				Method:     method,
				FullPath:   info.Path + "." + method.Name,
				PkgPath:    pkgPath,
				CallerPath: key,
			}

			callerCache.Store(key, callerEntry)
		}

	}

}
func GetAllCaller() []CallerEntry {
	ret := []CallerEntry{}
	callerCache.Range(func(key, value interface{}) bool {
		ret = append(ret, value.(CallerEntry))
		return true
	})
	return ret
}
