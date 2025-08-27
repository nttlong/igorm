package wx

// import (
// 	"fmt"
// 	"net/http"
// 	"reflect"
// 	"strings"
// 	"sync"
// 	"wx/handlers"
// )

// type Inject[TContainer any] struct {
// 	Context *HttpContext
// }

// var cachInjectNew map[reflect.Type]reflect.Value = map[reflect.Type]reflect.Value{}

// type initConatinerRegister struct {
// 	once sync.Once
// }

// var cacheInjectRegister sync.Map

// func (c *Inject[TContainer]) Register(fn func(injector *TContainer) error) {
// 	actually, _ := cacheInjectRegister.LoadOrStore(reflect.TypeFor[TContainer](), &initConatinerRegister{})
// 	actually.(*initConatinerRegister).once.Do(func() {
// 		cachInjectNew[reflect.TypeFor[TContainer]()] = reflect.ValueOf(fn)
// 	})

// }
// func (c *Inject[TContainer]) NewInstance() (*TContainer, error) {
// 	entry, ok := cachInjectNew[reflect.TypeFor[TContainer]()]
// 	if !ok {
// 		var ret TContainer
// 		return &ret, nil
// 	}
// 	arg := reflect.New(reflect.TypeFor[TContainer]())
// 	ret := entry.Call([]reflect.Value{arg})
// 	if ret[1].IsNil() {
// 		return ret[0].Interface().(*TContainer), nil
// 	}
// 	return nil, ret[1].Interface().(error)

// }

// var pkgPathOfInject = reflect.TypeOf(Inject[any]{}).PkgPath()
// var nameOfInject = strings.Split(reflect.TypeOf(Inject[any]{}).Name(), "[")[0] + "["

// type initIsInjectType struct {
// 	once sync.Once
// 	ret  bool
// }

// var cacheIsInjectType sync.Map

// func isInjectType(typ reflect.Type) bool {
// 	actually, _ := cacheIsInjectType.LoadOrStore(typ, &initIsInjectType{})
// 	oneItem := actually.(*initIsInjectType)
// 	oneItem.once.Do(func() {

// 		if typ.Kind() == reflect.Ptr {
// 			typ = typ.Elem()
// 		}
// 		for i := 0; i < typ.NumField(); i++ {
// 			field := typ.Field(i)
// 			if field.Anonymous {
// 				fieldType := field.Type
// 				if fieldType.Kind() == reflect.Ptr {
// 					fieldType = fieldType.Elem()
// 				}
// 				if fieldType.Kind() == reflect.Struct {
// 					if fieldType.PkgPath() == pkgPathOfInject && strings.HasPrefix(fieldType.Name(), nameOfInject) {
// 						oneItem.ret = true
// 						return
// 					}
// 				}

// 			}
// 		}
// 		oneItem.ret = false
// 	})
// 	return oneItem.ret
// }
// func isReadyRegister(typ reflect.Type) bool {
// 	if typ.Kind() == reflect.Ptr {
// 		typ = typ.Elem()
// 	}
// 	if _, ok := cachInjectNew[typ]; ok {
// 		return true
// 	}

// 	return false
// }
// func NewInjectByType(typ reflect.Type, r *http.Request, w http.ResponseWriter) (reflect.Value, error) {
// 	typEle := typ
// 	if typEle.Kind() == reflect.Ptr {
// 		typEle = typEle.Elem()
// 	}

// 	entry, ok := cachInjectNew[typEle]
// 	if !ok {

// 		return reflect.New(typ), fmt.Errorf("%s not register", typ.String())
// 	}
// 	arg := reflect.New(typEle)
// 	contextField := arg.Elem().FieldByName("Context")
// 	if contextField.IsValid() {
// 		contextField.Set(reflect.ValueOf(&HttpContext{
// 			Request:  r,
// 			Response: w,
// 		}))
// 	}
// 	ret := entry.Call([]reflect.Value{arg})
// 	if len(ret) == 1 && ret[0].IsNil() {
// 		if ret[0].IsNil() {
// 			return arg, nil
// 		} else if err, ok := ret[0].Interface().(error); ok {
// 			return arg, err
// 		}

// 	}
// 	return arg, nil

// }
// func init() {
// 	handlers.IsInjectType = isInjectType
// 	handlers.IsReadyRegister = isReadyRegister
// 	handlers.NewInjectByType = NewInjectByType

// }
