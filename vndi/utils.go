package vdi

import (
	"reflect"
)

func factoryType(factory any) reflect.Type {
	t := reflect.TypeOf(factory)
	if t.Kind() != reflect.Func || t.NumOut() != 1 {
		panic("invalid factory")
	}
	return t.Out(0)
}

func callFactoryMust(factory any, owner Container) any {
	fVal := reflect.ValueOf(factory)
	fType := fVal.Type()
	if fType.NumIn() == 0 {
		return fVal.Call(nil)[0].Interface()
	}
	return fVal.Call([]reflect.Value{reflect.ValueOf(owner)})[0].Interface()
}

func TypeOf[T any]() reflect.Type {
	return reflect.TypeFor[T]()
}
