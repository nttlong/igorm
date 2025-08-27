package wx

import (
	"fmt"
	"net/http"
	"reflect"
	"sync"
	"wx/handlers"
)

type Service struct {
	Req *http.Request
	Res http.ResponseWriter
}

type serviceHelperType struct {
}
type initServiceHelperTypeIsTypeHasServiceContext struct {
	init bool
	once sync.Once
}

var cacheServiceHelperTypeIsTypeHasServiceContext sync.Map

func (svcHelper *serviceHelperType) isTypeHasServiceContextInternal(typ reflect.Type) bool {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return false
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Anonymous {
			fieldType := field.Type
			if fieldType.Kind() == reflect.Ptr {
				fieldType = fieldType.Elem()
			}
			if fieldType == reflect.TypeOf(Service{}) {
				return true
			}
		}

	}
	return false
}
func (svcHelper *serviceHelperType) IsTypeHasServiceContext(typ reflect.Type) bool {
	actually, _ := cacheServiceHelperTypeIsTypeHasServiceContext.LoadOrStore(typ, &initServiceHelperTypeIsTypeHasServiceContext{})
	item := actually.(*initServiceHelperTypeIsTypeHasServiceContext)
	item.once.Do(func() {
		item.init = svcHelper.isTypeHasServiceContextInternal(typ)
	})
	return item.init
}

type initServiceHelperTypeFindNewMethod struct {
	init *reflect.Method
	err  error
	once sync.Once
}

var cacheServiceHelperTypeFindNewMethod sync.Map

func (svcHelper *serviceHelperType) FindNewMethod(typ reflect.Type) (*reflect.Method, error) {
	actually, _ := cacheServiceHelperTypeFindNewMethod.LoadOrStore(typ, &initServiceHelperTypeFindNewMethod{})
	item := actually.(*initServiceHelperTypeFindNewMethod)
	item.once.Do(func() {
		item.init, item.err = svcHelper.findNewMethodInternal(typ)
	})
	return item.init, item.err
}
func (svcHelper *serviceHelperType) findNewMethodInternal(typ reflect.Type) (*reflect.Method, error) {
	prtType := typ
	if prtType.Kind() == reflect.Struct {
		prtType = reflect.PointerTo(prtType)
	}
	for i := 0; i < prtType.NumMethod(); i++ {
		ret := prtType.Method(i)
		if ret.Name == "New" {
			if ret.Type.NumOut() != 1 {
				return nil, fmt.Errorf("%s.New must return error type only", prtType.String())
			}
			return &ret, nil

		}
	}
	return nil, fmt.Errorf("%s.New was not found", prtType.String())
}
func (svcHelper *serviceHelperType) NewInjectByType(typ reflect.Type, newMethod reflect.Method, r *http.Request, w http.ResponseWriter) (*reflect.Value, error) {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	ret := reflect.New(typ)
	reqField := ret.Elem().FieldByName("Req")
	if reqField.IsValid() {
		reqField.Set(reflect.ValueOf(r))
	}
	resField := ret.Elem().FieldByName("Res")
	if resField.IsValid() {
		resField.Set(reflect.ValueOf(w))
	}
	retCall := newMethod.Func.Call([]reflect.Value{ret})
	if retCall[0].IsNil() {
		return &ret, nil
	} else {
		return nil, retCall[0].Interface().(error)
	}

}

var serviceHelper = &serviceHelperType{}

func init() {
	handlers.IsServiceContext = serviceHelper.IsTypeHasServiceContext
	handlers.GetNewMethodOfServiceContext = serviceHelper.FindNewMethod
	handlers.NewInjectByType = serviceHelper.NewInjectByType

}
