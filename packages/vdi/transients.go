package vdi

import (
	"fmt"
	"reflect"
	"sync"
)

type BaseService struct {
	Err             error
	Transisent      map[reflect.Type]reflect.Value
	TransisentError map[reflect.Type]error
	SingletonFunc   map[reflect.Type]reflect.Value
	SingletonValue  map[reflect.Type]interface{}
}
type initTransisentRegister struct {
	once     sync.Once
	instance interface{}
	err      error
}

var cacheTransisentRegister sync.Map

type transisentRegisterKey struct {
	typ       reflect.Type
	container interface{}
}
type mapKeyTransisentRegisterStruct struct {
	typ       reflect.Type
	container interface{}
}

var mapKeyTransisentRegister map[reflect.Type]map[interface{}]mapKeyTransisentRegisterStruct

func TransisentRegister[T any](container interface{}, fn func() (*T, error)) {
	if mapKeyTransisentRegister == nil {
		mapKeyTransisentRegister = make(map[reflect.Type]map[interface{}]mapKeyTransisentRegisterStruct)
	}
	if _, ok := mapKeyTransisentRegister[reflect.TypeFor[T]()]; !ok {
		mapKeyTransisentRegister[reflect.TypeFor[T]()] = make(map[interface{}]mapKeyTransisentRegisterStruct)
	}
	mapKeyTransisentRegister[reflect.TypeFor[T]()][container] = mapKeyTransisentRegisterStruct{
		typ:       reflect.TypeFor[T](),
		container: container,
	}

	key := mapKeyTransisentRegister[reflect.TypeFor[T]()][container]
	actual, _ := cacheTransisentRegister.LoadOrStore(key, &initTransisentRegister{})
	initTransisent := actual.(*initTransisentRegister)
	initTransisent.once.Do(func() {
		TransisentField := reflect.ValueOf(container).Elem().FieldByName("Transisent")
		TransisentVal := TransisentField.Addr().Interface().(*map[reflect.Type]reflect.Value)
		if *TransisentVal == nil {
			*TransisentVal = make(map[reflect.Type]reflect.Value)
		}
		typ, fnVal := reflect.TypeFor[T](), reflect.ValueOf(fn)
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		(*TransisentVal)[typ] = fnVal

	})
	if initTransisent.err != nil {
		panic(initTransisent.err)
	}

}

func TransisentGet[T any](container interface{}) T {
	TransisentField := reflect.ValueOf(container).Elem().FieldByName("Transisent")
	TransisentVal := TransisentField.Addr().Interface().(*map[reflect.Type]reflect.Value)
	if *TransisentVal == nil {
		panic("container is nil")
	}
	typ := reflect.TypeFor[T]()
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	fn, ok := (*TransisentVal)[reflect.TypeFor[T]()]
	if !ok {
		typList := ""
		for k := range *TransisentVal {
			typList += k.String() + " \n"
		}
		panic(fmt.Errorf("container is nil, type %s not found in Transisent, available types are: \n\t%s", reflect.TypeFor[T]().String(), typList))
	}
	if fn.IsValid() {
		ret := fn.Call([]reflect.Value{})
		if len(ret) == 2 {
			if !ret[1].IsNil() {
				err := ret[1].Interface().(error)
				panic(err)

			}
			ret := ret[0].Interface().(*T)
			return *ret
		}

	}
	panic("container is nil")
}
func SingletonRegisterFunc[T any](container interface{}, fn func() *T) {
	SingletonField := reflect.ValueOf(container).Elem().FieldByName("SingletonFunc")
	SingletonVal := SingletonField.Addr().Interface().(*map[reflect.Type]reflect.Value)
	if *SingletonVal == nil {
		*SingletonVal = make(map[reflect.Type]reflect.Value)
	}

	(*SingletonVal)[reflect.TypeFor[T]()] = reflect.ValueOf(fn)

}

type initSingletonGet struct {
	once sync.Once
	val  interface{}
	err  error
}

var cacheSingletonGet sync.Map

func SingletonGet[T any](container interface{}) *T {
	actual, _ := cacheSingletonGet.LoadOrStore(reflect.TypeFor[T](), &initSingletonGet{})
	initSingleton := actual.(*initSingletonGet)
	initSingleton.once.Do(func() {
		initSingleton.val = singletonGet[T](container)
	})
	if initSingleton.err != nil {
		panic(initSingleton.err)
	}
	return initSingleton.val.(*T)

}
func singletonGet[T any](container interface{}) *T {
	SingletonField := reflect.ValueOf(container).Elem().FieldByName("SingletonFunc")
	SingletonVal := SingletonField.Addr().Interface().(*map[reflect.Type]reflect.Value)
	if *SingletonVal == nil {
		panic("container is nil")
	}
	fn := (*SingletonVal)[reflect.TypeFor[T]()]
	if fn.IsValid() {
		ret := fn.Call([]reflect.Value{})
		if len(ret) == 1 {
			ret := ret[0].Interface().(*T)
			return ret
		}

	}
	panic("container is nil")
}
