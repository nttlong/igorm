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

func GetInputTypeOfCallerPath(callerPath string) ([]reflect.Type, *reflect.Method, error) {
	if !strings.Contains(callerPath, "@") {
		return nil, nil, CallError{
			Err:  fmt.Errorf("caller path should be in format of package.path@method"),
			Code: CallErrorCodeInvalidCallerPath,
		}
	}
	callerEntry, found := callerCache.Load(strings.ToLower(callerPath))
	if !found {
		return nil, nil, CallError{
			Err:  fmt.Errorf("caller not found"),
			Code: CallErrorCodeCallerPathNotFound,
		}
	}
	method := callerEntry.(CallerEntry).Method
	ret := make([]reflect.Type, method.Type.NumIn()-1)
	for i := 1; i < method.Type.NumIn(); i++ {
		ret[i-1] = method.Type.In(i)
	}
	return ret, &method, nil
}

// type InjectorCaller struct {
// 	callerPath string
// 	injector   interface{}
// }

func NewCaller(callerPath string, injector interface{}) func(args interface{}) (interface{}, error) {
	return func(args interface{}) (interface{}, error) {
		return Call(callerPath, args, injector)
	}
}

//	func (i *InjectorCaller) Call(args interface{}) (interface{}, error) {
//		return Call(i.callerPath, args, i.injector)
//	}
func Call(callerPath string, args interface{}, injector interface{}) (interface{}, error) {

	// if req, ok := args.(*RequestType); ok {
	// 	args = req.Get("Args")
	// } else {
	// 	argsType := reflect.TypeOf(args)
	// 	if argsType.Kind() == reflect.Ptr {
	// 		argsType = argsType.Elem()
	// 	if argsType.Kind() != reflect.Slice {
	// 		argsValue := reflect.ValueOf(args)

	// 		if argsValue.Kind() == reflect.Ptr {
	// 			argsValue = argsValue.Elem()
	// 			argsType = argsType.Elem()

	// 		}
	// 		if field, found := argsType.FieldByName("Args"); found {
	// 			argsValue = argsValue.FieldByName(field.Name)
	// 			if argsValue.Kind() == reflect.Ptr {
	// 				argsValue = argsValue.Elem()
	// 			}
	// 			args = argsValue.Interface()
	// 		}
	// 	}

	// 	if !strings.Contains(callerPath, "@") {
	// 		return nil, CallError{
	// 			Err:  fmt.Errorf("caller path should be in format of package.path@method"),
	// 			Code: CallErrorCodeInvalidCallerPath,
	// 		}
	// 	}
	// }
	fx := reflect.ValueOf(args)
	if fx.Kind() == reflect.Ptr {
		fx = fx.Elem()
		args = fx.Interface()

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

func isMethodFromStructOnly(t reflect.Type, methodName string) bool {
	// Kiểm tra xem type có phải struct không
	t1 := t
	if t.Kind() != reflect.Struct {
		t1 = t.Elem()
	}

	// Duyệt qua tất cả field để tìm embed
	for i := 0; i < t1.NumField(); i++ {
		field := t1.Field(i)
		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			// Kiểm tra method của embed hiện tại
			embedType := field.Type
			for j := 0; j < embedType.NumMethod(); j++ {
				embedMethod := embedType.Method(j)
				if embedMethod.Name == methodName {
					return false // Method thuộc về embed
				}
			}
			if field.Type.Kind() == reflect.Ptr {
				embedType = field.Type.Elem()
			}
			for j := 0; j < embedType.NumMethod(); j++ {
				embedMethod := embedType.Method(j)
				if embedMethod.Name == methodName {
					return false // Method thuộc về embed
				}
			}
			ptrEmbedInstance := reflect.New(embedType)
			ptrEmbedTyp := ptrEmbedInstance.Type()

			for j := 0; j < ptrEmbedTyp.NumMethod(); j++ {
				embedMethod := ptrEmbedTyp.Method(j)
				if embedMethod.Name == methodName {
					return false // Method thuộc về embed
				}
			}
			// Đệ quy kiểm tra embed lồng nhau
			if !isMethodFromStructOnly(embedType, methodName) {
				return false // Method thuộc về embed lồng nhau
			}
		}
	}

	// Nếu không tìm thấy method trong bất kỳ embed nào, trả về true
	return true
}
func getAllMethods(caller interface{}) []reflect.Method {
	ret := []reflect.Method{}
	typ := reflect.TypeOf(caller)

	fmt.Print(typ.Name())
	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)

		if isMethodFromStructOnly(typ, method.Name) {
			//check if first letter is uppercase
			if method.Name[0] >= 'A' && method.Name[0] <= 'Z' {
				ret = append(ret, method)
			}
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
