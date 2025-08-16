package wx

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type Depend[TInstance any, TApp any] struct {
	ins    TInstance
	err    error
	App    interface{}
	fnInit func(app *TApp) TInstance
	once   sync.Once
}

func (d *Depend[TInstance, TApp]) Ins() (*TInstance, error) {
	d.once.Do(func() {
		if d.fnInit == nil {
			t1 := reflect.TypeFor[TInstance]()
			if t1.Kind() == reflect.Ptr {
				t1 = t1.Elem()
			}
			t2 := reflect.TypeFor[TApp]()
			if t2.Kind() == reflect.Ptr {
				t2 = t2.Elem()
			}

			d.err = fmt.Errorf("please call Init of %s when New of %s is called", t1.String(), t2.String())
			return
		}
		d.ins = d.fnInit(d.App.(*TApp))
	})
	return &d.ins, d.err
}
func (d *Depend[TInstance, TApp]) Init(fn func(app *TApp) TInstance) {
	d.fnInit = fn
}
func findNewMethod[T any]() (*reflect.Method, error) {
	t := reflect.TypeFor[*T]()
	fmt.Println(t.String())
	for i := 0; i < t.NumMethod(); i++ {
		if t.Method(i).Name == "New" {
			ret := t.Method(i)
			if ret.Type.NumOut() != 1 {
				return nil, fmt.Errorf("function New of %s must return error", t.String())
			}
			return &ret, nil
		}
	}
	return nil, fmt.Errorf("%s of %T was not found", "New", *new(T))
}
func isTypeDepen(t reflect.Type, visited map[reflect.Type]bool) bool {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	visited[t] = true
	check := t.Kind() == reflect.Struct && t.PkgPath() == reflect.TypeOf(Depend[any, any]{}).PkgPath() && strings.HasPrefix(t.Name(), "Depend[")
	if check {
		return true
	} else {
		for i := 0; i < t.NumField(); i++ {
			if isTypeDepen(t.Field(i).Type, visited) {
				return true
			}
		}
		return false
	}
}
func createDepen(appInstance reflect.Value, depenType reflect.Type) reflect.Value {
	if depenType.Kind() == reflect.Ptr {
		depenType = depenType.Elem()
	}
	ret := reflect.New(depenType)
	fieldSet := ret.Elem().FieldByName("App")
	if fieldSet.IsValid() {
		fieldSet.Set(appInstance)
	} else {
		fieldSet.Set(reflect.ValueOf(appInstance.Interface()))

	}
	// ret.Elem().FieldByName("App").Set(appInstance)
	return ret
}
func Start[TApp any](fn func(app *TApp) error) error {
	newMethod, err := findNewMethod[TApp]()
	fmt.Println(newMethod)
	if err != nil {
		return err
	}
	typ := reflect.TypeFor[TApp]()
	appInstance := reflect.New(typ)

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if !isTypeDepen(field.Type, map[reflect.Type]bool{}) {
			continue
		}
		instanceVal := createDepen(appInstance, field.Type)
		fieldSet := appInstance.Elem().FieldByName(field.Name)
		if fieldSet.IsValid() {
			if instanceVal.Kind() == reflect.Ptr {
				fieldSet.Set(instanceVal)
			} else {
				fieldSet.Set(instanceVal.Elem())
			}

		} else {
			fieldSet.Set(reflect.ValueOf(instanceVal.Interface()))
		}

	}
	retCall := newMethod.Func.Call([]reflect.Value{appInstance})
	if len(retCall) == 1 && !retCall[0].IsNil() {
		return retCall[0].Interface().(error)
	}
	err = fn(appInstance.Interface().(*TApp))
	if err != nil {
		return err
	}

	return nil
}
