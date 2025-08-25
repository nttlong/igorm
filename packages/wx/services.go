package wx

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"wx/handlers"
)

type HttpContext struct {
	Request  *http.Request
	Response http.ResponseWriter
}
type HttpService[TService any] struct {
	instance    *TService
	Err         error
	newMethod   *reflect.Method // The method to create a new instance of the service
	HttpContext *HttpContext
}

func findNewMethodOfHttpServiceInternal(serviceType reflect.Type) (*reflect.Method, error) {
	t := serviceType
	if t.Kind() == reflect.Struct {
		t = reflect.PointerTo(t)
	}
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		if method.Name == "New" {
			return &method, nil
		}
	}
	return nil, nil

}

type initFindNewMethodOfHttpService struct {
	val  *reflect.Method
	err  error
	once sync.Once
}

var cachedFindNewMethodOfHttpService sync.Map

func findNewMethodOfHttpService(serviceType reflect.Type) (*reflect.Method, error) {
	typ := serviceType
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	actual, _ := cachedFindNewMethodOfHttpService.LoadOrStore(typ, &initFindNewMethodOfHttpService{})
	init := actual.(*initFindNewMethodOfHttpService)
	init.once.Do(func() {
		method, err := findNewMethodOfHttpServiceInternal(serviceType)
		if err != nil {
			init.err = err
			return
		}
		init.val = method
	})
	return init.val, init.err
}

func (hs *HttpService[TService]) Ins() (*TService, error) {
	if hs.instance == nil {
		if hs.newMethod == nil {
			newMethod, err := findNewMethodOfHttpService(reflect.TypeFor[TService]())

			if err != nil {
				return nil, err
			}
			if newMethod == nil {
				return nil, errors.New("method New was not found in inject HttpService[" + reflect.TypeFor[TService]().Name() + "]") // No New method found
			}
			hs.newMethod = newMethod
		}

		info, err := ServiceUtil.ExtractInfo(*hs.newMethod)
		if err != nil {
			return nil, err
		}
		args := make([]reflect.Value, (*hs.newMethod).Type.NumIn())
		reciverType := hs.newMethod.Type.In(0)
		if reciverType.Kind() == reflect.Ptr {
			reciverType = reciverType.Elem()
		}
		args[0] = reflect.New(reciverType) // The first argument is the receiver (the service type itself)
		if (*hs.newMethod).Type.In(info.IndexOfHttpContext).Kind() == reflect.Ptr {

			args[info.IndexOfHttpContext] = reflect.ValueOf(hs.HttpContext) // Create a new HttpContext pointer
		} else {
			args[info.IndexOfHttpContext] = reflect.ValueOf(*hs.HttpContext) // Create a new HttpContext value
		}

		// Set the HttpContext argument
		for _, index := range info.IndexOfInjectords {
			x, err := Helper.DependNew((*hs.newMethod).Type.In(index))
			if err != nil {
				return nil, fmt.Errorf("error creating dependency for %s: %w", (*hs.newMethod).Type.In(index).String(), err)
			}
			args[index] = *x
		}
		//"reflect: Call using **mockcontroller.UserService as type *mockcontroller.UserService"

		rets := (*hs.newMethod).Func.Call(args)
		if len(rets) == 1 {
			if err, ok := rets[0].Interface().(error); ok {
				return nil, err
			}
		}

		return args[0].Interface().(*TService), nil

	}
	return hs.instance, hs.Err
}

var httpServiceName = strings.Split(reflect.TypeFor[HttpService[any]]().Name(), "[")[0] + "["
var httpServicePath = reflect.TypeFor[HttpService[any]]().PkgPath()

func createServiceContext(req *http.Request, res http.ResponseWriter) reflect.Value {
	ret := &HttpContext{
		Request:  req,
		Response: res,
	}
	return reflect.ValueOf(ret)
}
func init() {
	handlers.DependIsHttpService = ServiceUtil.isInjectHttpService
	handlers.CreateServiceContext = createServiceContext
	handlers.DependIsHttpServiceMethodHasContext = dependIsHttpServiceMethodHasContext
}
