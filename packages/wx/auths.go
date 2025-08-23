package wx

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"wx/handlers"
	//_ "wx/handlers"
)

type AuthContext struct {
	Req *http.Request
	Res http.ResponseWriter
}
type Auth[T any] struct {
	fnInit   func(ctx *AuthContext) (*T, error)
	newMthod *reflect.Method
	Context  *AuthContext
}
type initAuthFindNewMehod struct {
	val  *reflect.Method
	err  error
	once sync.Once
}

var cacheAuthFindNewMehod sync.Map

func authFindNewMehod(typ reflect.Type) (*reflect.Method, error) {
	actually, _ := cacheAuthFindNewMehod.LoadOrStore(typ, &initAuthFindNewMehod{})
	init := actually.(*initAuthFindNewMehod)
	init.once.Do(func() {
		if typ.Kind() != reflect.Ptr {
			typ = reflect.PointerTo(typ)
		}
		for i := 0; i < typ.NumMethod(); i++ {
			if typ.Method(i).Name == "New" {
				ret := typ.Method(i)
				if ret.Type.NumIn() != 2 || ret.Type.NumOut() != 1 {
					init.err = fmt.Errorf("New method of %s must have 2 input parameters and 1 output parameters (return error)", typ.String())
					return
				}
				argsHandler := ret.Type.In(1)
				if argsHandler != reflect.TypeOf(&AuthContext{}) {
					init.err = fmt.Errorf("New method of %s must have AuthContext as second input parameter", typ.String())
					return
				}

				init.val = &ret
				return
			}
		}
		init.val = nil
	})
	return init.val, init.err
}

func (a *Auth[T]) Get() (*T, error) {
	if a.fnInit == nil {
		mt, err := authFindNewMehod(reflect.TypeFor[T]())
		if err != nil {
			return nil, err

		}
		retInsatnce, err := handlers.Helper.DepenAuthCreateInstance(*mt, reflect.ValueOf(a.Context))
		if err != nil {
			return nil, err
		}
		return retInsatnce.Interface().(*T), nil

	}
	return a.fnInit(a.Context)
}
func (a *Auth[T]) Init(fn func(ctx *AuthContext) (*T, error)) {
	a.fnInit = fn
}

var authPrefix = strings.Split(reflect.TypeOf(Auth[any]{}).Name(), "[")[0]
var authPackagePath = reflect.TypeOf(Auth[any]{}).PkgPath()

func depenAuthCreateCreateAuthContext(req *http.Request, res http.ResponseWriter) (*reflect.Value, error) {
	ctx := &AuthContext{
		Req: req,
		Res: res,
	}
	ret := reflect.ValueOf(ctx)
	return &ret, nil
}

func IsGeneriAuth(typ reflect.Type) bool {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return false
	}
	if typ.PkgPath() == authPackagePath && strings.HasPrefix(typ.Name(), authPrefix) {
		return true
	}
	return false

}
func init() {
	handlers.IsGeneriAuth = IsGeneriAuth
	handlers.DepenAuthCreateCreateAuthContext = depenAuthCreateCreateAuthContext
}
