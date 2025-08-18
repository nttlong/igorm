package handlers

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type initDependFIndNewMethod struct {
	val  *reflect.Method
	err  error
	once sync.Once
}

var cacheDependFIndNewMethod sync.Map

func (h *helperType) DependFindNewMethod(typ reflect.Type) (*reflect.Method, error) {
	key := typ
	if key.Kind() == reflect.Ptr {
		key = key.Elem()
	}
	actual, _ := cacheDependFIndNewMethod.LoadOrStore(key, &initDependFIndNewMethod{})
	init := actual.(*initDependFIndNewMethod)
	init.once.Do(func() {
		init.val, init.err = h.dependFIndNewMethod(typ)
	})
	return init.val, init.err
}
func (h *helperType) dependFIndNewMethod(typ reflect.Type) (*reflect.Method, error) {
	if typ.Kind() != reflect.Ptr {
		typ = reflect.PointerTo(typ)
	}
	for i := 0; i < typ.NumMethod(); i++ {
		if typ.Method(i).Name == "New" {
			mt := typ.Method(i)
			if mt.Type.NumOut() != 1 {
				msgErr := fmt.Sprintf("%s.New must return error", typ.String())
				return nil, errors.New(msgErr)
			}

			if mt.Type.Out(0) != h.ErrorType {
				msgErr := fmt.Sprintf("%s.New return %s, expected %s.New to return error", typ.String(), mt.Type.Out(0).String(), typ.String())
				return nil, errors.New(msgErr)
			}
			for j := 1; j < mt.Type.NumIn(); j++ {
				if !h.IsGenericDepen(mt.Type.In(j)) {
					msgErr := fmt.Sprintf("arg %d of %s.New must be generic type wx.Depend, not %s", j, typ.String(), mt.Type.In(j).String())
					return nil, errors.New(msgErr)
				}
			}

			return &mt, nil
		}
	}
	return nil, nil

}
func (h *helperType) LoadAllFieldsInternal(insVal reflect.Value) error {
	return h.depenLoadAllFieldsInternal(insVal, map[reflect.Type]bool{})
}
func (h *helperType) depenLoadAllFieldsInternal(insVal reflect.Value, visited map[reflect.Type]bool) error {

	for i := 0; i < insVal.Elem().NumField(); i++ {
		field := insVal.Elem().Field(i)
		fieldType := field.Type()
		checKType := fieldType
		if checKType.Kind() == reflect.Ptr {
			checKType = checKType.Elem()
		}
		if _, ok := visited[checKType]; ok {
			msg := "circulation reference:"
			for k := range visited {
				msg += k.String() + "-->"
			}
			msg += checKType.String()
			return errors.New(msg)

		}
		if h.IsGenericDepen(fieldType) {
			depenVal, err := h.dependNewInternal(fieldType, visited)
			if err != nil {
				return err
			}
			if fieldType.Kind() == reflect.Ptr {
				field.Set(*depenVal)
			} else {

				field.Set((*depenVal).Elem())
			}
		}
	}
	return nil
}

type initDependNewOnce struct {
	val  *reflect.Value
	err  error
	once sync.Once
}

var cacheDependNewOnce sync.Map

func (h *helperType) DependNewOnce(typ reflect.Type) (*reflect.Value, error) {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	actual, _ := cacheDependNewOnce.LoadOrStore(typ, &initDependNewOnce{})
	init := actual.(*initDependNewOnce)
	init.once.Do(func() {
		init.val, init.err = h.DependNew(typ)
	})
	return init.val, init.err

}
func (h *helperType) DependNew(typ reflect.Type) (*reflect.Value, error) {
	return h.dependNewInternal(typ, map[reflect.Type]bool{})
}
func (h *helperType) dependNewInternal(typ reflect.Type, visited map[reflect.Type]bool) (*reflect.Value, error) {
	checKType := typ
	if checKType.Kind() == reflect.Ptr {
		checKType = checKType.Elem()
	}
	if _, ok := visited[checKType]; ok {
		msg := "circulation reference:"
		for k := range visited {
			msg += k.String() + "-->"
		}
		msg += checKType.String()
		return nil, errors.New(msg)

	}
	visited[checKType] = true

	newMethod, err := h.DependFindNewMethod(typ)
	if err != nil {
		return nil, err
	}
	if typ.Kind() != reflect.Ptr {
		typ = reflect.PointerTo(typ)

	}
	ret := reflect.New(typ.Elem())
	if newMethod == nil {

		err := h.depenLoadAllFieldsInternal(ret, visited)
		if err != nil {
			return nil, err
		}
		return &ret, nil

	} else {

		err := h.depenLoadAllFieldsInternal(ret, visited)
		if err != nil {
			return nil, err
		}
		/*
		 Scan all arg
		*/
		args := make([]reflect.Value, newMethod.Type.NumIn())
		args[0] = ret
		for i := 1; i < newMethod.Type.NumIn(); i++ {
			depenVal, err := h.dependNewInternal(newMethod.Type.In(i), visited)
			if err != nil {
				return nil, err
			}
			args[i] = *depenVal
		}
		retCall := newMethod.Func.Call(args)
		if !retCall[0].IsNil() {
			return nil, retCall[0].Interface().(error)
		}
		return &ret, nil

	}

}
