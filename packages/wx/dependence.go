package wx

import (
	"reflect"
	"strings"
	"sync"
	httpServer "wx/HtttpServer"
	"wx/handlers"
)

type Depend[TInstance any, TOwner any] struct {
	ins    *TInstance
	err    error
	Owner  interface{}
	fnInit func(Owner *TOwner) (*TInstance, error)
	once   sync.Once
}

func (d *Depend[TInstance, TOwner]) GetOwner() *TOwner {
	return d.Owner.(*TOwner)
}
func (d *Depend[TInstance, TOwner]) Ins() (*TInstance, error) {
	d.once.Do(func() {
		if d.fnInit == nil {
			ret, err := handlers.Helper.DependNew(reflect.TypeFor[TInstance]())
			d.err = err
			d.ins = (*ret).Interface().(*TInstance)
			return
		}
		d.ins, d.err = d.fnInit(d.Owner.(*TOwner))
	})
	return d.ins, d.err
}
func (d *Depend[TInstance, TOwner]) Init(fn func(apOwnerp *TOwner) (*TInstance, error)) {

	d.fnInit = fn
}

func isStructDepenType(t reflect.Type) bool {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return false
	}
	checkIsDepend := t.Kind() == reflect.Struct && t.PkgPath() == reflect.TypeOf(Depend[any, any]{}).PkgPath() && strings.HasPrefix(t.Name(), "Depend[")
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
func NewDepen[T any]() (*T, error) {
	ret, err := handlers.Helper.DependNew(reflect.TypeFor[T]())
	if err != nil {
		return nil, err
	}
	return ret.Interface().(*T), nil

}
