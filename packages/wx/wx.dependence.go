package wx

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	httpServer "wx/HtttpServer"
	"wx/handlers"
)

type Depend[TInstance any] struct {
	ins    *TInstance
	err    error
	Owner  interface{}
	fnInit func() (*TInstance, error)
	once   sync.Once
}

func (d *Depend[TInstance]) Ins() (*TInstance, error) {
	d.once.Do(func() {
		if d.fnInit == nil {
			// detect if TInstance is function type
			if reflect.TypeFor[TInstance]().Kind() == reflect.Func {
				fnPtrType := reflect.TypeFor[*TInstance]()
				mt, ok := fnPtrType.MethodByName("New")
				if !ok {
					d.err = fmt.Errorf("cannot find New method for function type %s", reflect.TypeFor[TInstance]().String())
					return
				}
				ret := reflect.New(fnPtrType.Elem())
				mt.Func.Call([]reflect.Value{ret})
				d.ins = ret.Interface().(*TInstance)
				return

			}

			ret, err := handlers.Helper.DependNew(reflect.TypeFor[TInstance]())
			if err != nil {
				d.err = err
				return
			} else {
				if ret == nil {
					d.err = fmt.Errorf("cannot create instance of type %s, method New was not found", reflect.TypeFor[TInstance]().String())
					return
				}
			}
			d.ins = (*ret).Interface().(*TInstance)
			return
		}
		d.ins, d.err = d.fnInit()
	})
	return d.ins, d.err
}
func (d *Depend[TInstance]) Init(fn func() (*TInstance, error)) {

	d.fnInit = fn
}

func isStructDepenType(t reflect.Type) bool {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return false
	}
	checkIsDepend := t.Kind() == reflect.Struct && t.PkgPath() == reflect.TypeOf(Depend[any]{}).PkgPath() && strings.HasPrefix(t.Name(), "Depend[")
	checkIsGlobal := t.Kind() == reflect.Struct && t.PkgPath() == reflect.TypeOf(Global[any]{}).PkgPath() && strings.HasPrefix(t.Name(), "Global[")
	return checkIsDepend || checkIsGlobal
}

func isTypeDepen(t reflect.Type, visited map[reflect.Type]bool) bool {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return false
	}
	visited[t] = true

	if isStructDepenType(t) {
		return true
	} else {
		for i := 0; i < t.NumField(); i++ {

			ft := t.Field(i).Type
			if ft.Kind() == reflect.Ptr {
				ft = ft.Elem()
			}
			if visited[ft] {
				continue
			}
			if ft.Kind() != reflect.Struct {
				continue
			}
			if isTypeDepen(t.Field(i).Type, visited) {
				return true
			}
		}
		return false
	}
}
func findNewMethodOfType(t reflect.Type) *reflect.Method {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil
	}
	prtT := reflect.PointerTo(t)
	for i := 0; i < prtT.NumMethod(); i++ {
		if prtT.Method(i).Name == "New" {
			ret := prtT.Method(i)
			return &ret
		}
	}
	return nil
}
func runNewMethod(ins reflect.Value, newMethod reflect.Method) error {
	args := make([]reflect.Value, newMethod.Type.NumIn())
	args[0] = ins
	for i := 0; i < newMethod.Type.NumIn(); i++ {
		nmt := findNewMethodOfType(newMethod.Type.In(i))
		if nmt != nil {
			err := runNewMethod(ins, *nmt)
			if err != nil {
				return err
			}
		} else {
			args[i] = reflect.New(newMethod.Type.In(i))
		}
	}
	return nil
}

func createDepen(depenType reflect.Type) (*reflect.Value, error) {
	if depenType.Kind() == reflect.Ptr {
		depenType = depenType.Elem()
	}

	ret, err := DepenResolvers.ResolveTypeOnce(depenType)
	if err != nil {
		return nil, err
	}

	return ret, nil
}
func Start[TApp any](fn func(app *TApp) error) error {
	ret, err := New[TApp]()
	if err != nil {
		return err
	}
	err = fn(ret)
	return nil

}
func init() {
	httpServer.IsTypeDepen = isTypeDepen
	httpServer.CreateDepen = createDepen
	httpServer.FindNewMethod = DepenResolvers.FindNewMethod
	httpServer.ResolveNewMethod = DepenResolvers.RunNewMethod
	httpServer.ResolveNewMethodWithReceiver = DepenResolvers.RunNewMethodWithReceiver
	handlers.IsGenericDepen = func(typ reflect.Type) bool {
		return isStructDepenType(typ)
	}

}
func NewDepend[T any]() (*T, error) {
	ret, err := handlers.Helper.DependNew(reflect.TypeFor[T]())
	if err != nil {
		return nil, err
	}
	return ret.Interface().(*T), nil

}
func NewGlobal[T any]() (*T, error) {
	ret, err := Helper.DependNewOnce(reflect.TypeFor[*T]())
	if err != nil {
		return nil, err
	}
	return ret.Interface().(*T), nil
}
