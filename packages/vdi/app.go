package vdi

import (
	"fmt"
	"reflect"
)

func getNewMethod[T any]() (*reflect.Method, error) {
	typ := reflect.TypeFor[*T]()
	var newMethod *reflect.Method
	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)
		if method.Name == "New" {
			newMethod = &method
		}

	}
	if newMethod == nil {
		return nil, fmt.Errorf("no new method found in %s", typ.String())
	}
	if newMethod.Type.NumIn() != 1 {
		return nil, fmt.Errorf("new method must have one argument. That arg is %s", typ.String())
	}
	return newMethod, nil

}
func Start[T any]() (*T, error) {
	newMethod, err := getNewMethod[T]()
	if err != nil {
		return nil, err
	}

	c := &Container[T]{}
	key := reflect.TypeFor[T]().String()
	actual, _ := initRegisterContainerCache.LoadOrStore(key, &initRegisterContainer{})
	initContainer := actual.(*initRegisterContainer)
	wrappperResolver := func(svc *T) error {
		containerField := reflect.ValueOf(svc).Elem().FieldByName("Container")
		if containerField.IsValid() {
			containerField.Set(reflect.ValueOf(c))
		}
		ret := newMethod.Func.Call([]reflect.Value{reflect.ValueOf(svc)})
		if len(ret) != 1 {
			return fmt.Errorf("new method must return one value. That value is %s", key)
		}
		if !ret[0].IsNil() {
			return ret[0].Interface().(error)
		}
		return nil
	}
	initContainer.once.Do(func() {
		initContainer.instance, initContainer.err = registerContainer(wrappperResolver)
	})

	if initContainer.err != nil {
		// ret := reflect.New(reflect.TypeFor[T]()).Interface()
		// f := reflect.ValueOf(ret).Elem().FieldByName("Error")
		// f.Set(reflect.ValueOf(initContainer.err))
		// tRet := ret.(*T)

		return nil, initContainer.err

	} else {
		innstanceField := reflect.ValueOf(initContainer.instance).FieldByName("Instance")
		ret := innstanceField.Elem().Interface().(*T)

		return ret, nil
	}

}
