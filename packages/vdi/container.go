package vdi

import (
	"fmt"
	"reflect"
	"sync"
)

type Container[T any] struct {
	Error       error
	Transitions map[reflect.Type]reflect.Value
}

func (c *Container[T]) GetTransition(typ reflect.Type) (interface{}, error) {
	if c == nil {
		return nil, fmt.Errorf("container is nil")
	}
	if c.Error != nil {
		return nil, c.Error
	}
	if fn, ok := c.Transitions[typ]; ok {
		return fn.Call([]reflect.Value{})[0].Interface(), nil
	}
	return nil, fmt.Errorf("no transition found for type %s in %s", typ.String(), reflect.TypeOf(c).Elem().String())

}
func GetTransitionFromContainer[T any](container interface{}) (*T, error) {
	containerVal := reflect.ValueOf(container)
	if containerVal.Kind() == reflect.Ptr {
		containerVal = containerVal.Elem()
	}

	TransitionsField := containerVal.FieldByName("Transitions")
	if !TransitionsField.IsValid() {
		panic(fmt.Errorf("param container is not a valid container, pls check the type of container %T", reflect.TypeOf(container)))
	}
	mapInstance := TransitionsField.Interface().(map[reflect.Type]reflect.Value)
	outputType := reflect.TypeFor[T]()
	if outputType.Kind() == reflect.Ptr {
		outputType = outputType.Elem()
	}
	if fnVal, ok := mapInstance[outputType]; ok {
		return fnVal.Call([]reflect.Value{})[0].Interface().(*T), nil
	} else {
		return nil, fmt.Errorf("no transition found for type %s in %s", outputType.String(), reflect.TypeOf(container).Elem().String())
	}

}

type initAddTransitionToContainer struct {
	once     sync.Once
	instance interface{}
	err      error
}

var cacheAddTransitionToContainer sync.Map

func AddTransitionToContainer[T any](container interface{}, fn func() (*T, error)) interface{} {
	actual, _ := cacheAddTransitionToContainer.LoadOrStore(reflect.TypeFor[T](), &initAddTransitionToContainer{})
	initContainer := actual.(*initAddTransitionToContainer)
	initContainer.once.Do(func() {
		initContainer.instance = addTransitionToContainer(container, fn)
	})
	return initContainer.instance

}
func addTransitionToContainer[T any](container interface{}, fn func() (*T, error)) interface{} {
	containerVal := reflect.ValueOf(container)
	if containerVal.Kind() == reflect.Ptr {
		containerVal = containerVal.Elem()
	}

	TransitionsField := containerVal.FieldByName("Transitions")
	if !TransitionsField.IsValid() {
		panic(fmt.Errorf("param container is not a valid container, pls check the type of container %T", reflect.TypeOf(container)))

	}
	fmt.Println(TransitionsField.IsValid())
	mapInstance := TransitionsField.Interface().(map[reflect.Type]reflect.Value)
	fmt.Println(mapInstance)
	val := reflect.ValueOf(fn)
	outputType := val.Type().Out(0)
	if outputType.Kind() == reflect.Ptr {
		outputType = outputType.Elem()
	}
	fmt.Print(outputType.String())
	mapInstance[outputType] = val
	return container

}

var initRegisterContainerCache = sync.Map{}

func (c *Container[T]) New(resolver func(owner *T) error) *T {
	// if c == nil {
	// 	fmt.Printf("container is nil ,Container[%s]", reflect.TypeFor[T]().String())
	// 	return nil
	// }

	key := reflect.TypeFor[T]().String()
	actual, _ := initRegisterContainerCache.LoadOrStore(key, &initRegisterContainer{})
	initContainer := actual.(*initRegisterContainer)
	wrappperResolver := func(svc *T) error {
		containerField := reflect.ValueOf(svc).Elem().FieldByName("Container")
		if !containerField.IsValid() {
			panic(fmt.Errorf("param svc is not a valid svc, pls check the type of svc %T", reflect.TypeOf(svc)))
		}
		containerField.Set(reflect.ValueOf(c))
		return resolver(svc)
	}
	initContainer.once.Do(func() {
		initContainer.instance, initContainer.err = registerContainer(wrappperResolver)
	})

	if initContainer.err != nil {
		ret := reflect.New(reflect.TypeFor[T]()).Interface()
		f := reflect.ValueOf(ret).Elem().FieldByName("Error")
		f.Set(reflect.ValueOf(initContainer.err))
		tRet := ret.(*T)

		return tRet

	} else {
		innstanceField := reflect.ValueOf(initContainer.instance).FieldByName("Instance")
		ret := innstanceField.Elem().Interface().(*T)

		// f := reflect.ValueOf(ret).Elem().FieldByName("Container")
		// if f.Kind() == reflect.Ptr {
		// 	f = f.Elem()
		// }
		// cVal := reflect.ValueOf(c)
		// if cVal.Kind() == reflect.Ptr {
		// 	cVal = cVal.Elem()
		// }
		// if cVal.CanAddr() {
		// 	if f.IsValid() && f.CanSet() {
		// 		f.Set(reflect.ValueOf(c).Elem())
		// 	} else {
		// 		panic(fmt.Errorf("field Container not found in type %s", reflect.TypeFor[T]().String()))
		// 	}
		// }
		return ret
	}

}
func NewContainer[T any](resolver func(owner *T) error) *T {
	// if c == nil {
	// 	fmt.Printf("container is nil ,Container[%s]", reflect.TypeFor[T]().String())
	// 	return nil
	// }
	c := &Container[T]{}
	key := reflect.TypeFor[T]().String()
	actual, _ := initRegisterContainerCache.LoadOrStore(key, &initRegisterContainer{})
	initContainer := actual.(*initRegisterContainer)
	wrappperResolver := func(svc *T) error {
		containerField := reflect.ValueOf(svc).Elem().FieldByName("Container")
		if containerField.IsValid() {
			containerField.Set(reflect.ValueOf(c))
		}

		return resolver(svc)
	}
	initContainer.once.Do(func() {
		initContainer.instance, initContainer.err = registerContainer(wrappperResolver)
	})

	if initContainer.err != nil {
		ret := reflect.New(reflect.TypeFor[T]()).Interface()
		f := reflect.ValueOf(ret).Elem().FieldByName("Error")
		f.Set(reflect.ValueOf(initContainer.err))
		tRet := ret.(*T)

		return tRet

	} else {
		innstanceField := reflect.ValueOf(initContainer.instance).FieldByName("Instance")
		ret := innstanceField.Elem().Interface().(*T)

		// f := reflect.ValueOf(ret).Elem().FieldByName("Container")
		// if f.Kind() == reflect.Ptr {
		// 	f = f.Elem()
		// }
		// cVal := reflect.ValueOf(c)
		// if cVal.Kind() == reflect.Ptr {
		// 	cVal = cVal.Elem()
		// }
		// if cVal.CanAddr() {
		// 	if f.IsValid() && f.CanSet() {
		// 		f.Set(reflect.ValueOf(c).Elem())
		// 	} else {
		// 		panic(fmt.Errorf("field Container not found in type %s", reflect.TypeFor[T]().String()))
		// 	}
		// }
		return ret
	}

}
