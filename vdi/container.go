package vdi

import (
	"reflect"
)

type Container[T any] struct {
	Error error
}

func (c *Container[T]) New(resolver func(svc *T) error) *T {
	key := reflect.TypeFor[T]().String()
	actual, _ := initRegisterContainerCache.LoadOrStore(key, &initRegisterContainer{})
	initContainer := actual.(*initRegisterContainer)
	initContainer.once.Do(func() {
		initContainer.instance, initContainer.err = registerContainer(resolver)
	})

	if initContainer.err != nil {
		ret := reflect.New(reflect.TypeFor[T]()).Interface()
		f := reflect.ValueOf(ret).Elem().FieldByName("Error")
		f.Set(reflect.ValueOf(initContainer.err))
		tRet := ret.(*T)

		return tRet

	} else {
		ret := initContainer.instance.(*containerInfo[T]).Get()
		return ret
	}

}
