package wx

import (
	"reflect"
	"sync"
)

type initGetMethod struct {
	once   sync.Once
	method *reflect.Method
}

var cacheGetMethod sync.Map

func GetMethod[T any](methodName string) *reflect.Method {
	t := reflect.TypeFor[*T]()
	typ := reflect.TypeFor[T]()
	atual, _ := cacheGetMethod.LoadOrStore(typ, &initGetMethod{})
	init := atual.(*initGetMethod)
	init.once.Do(func() {
		for i := 0; i < t.NumMethod(); i++ {
			if t.Method(i).Name == methodName {
				m := t.Method(i)
				init.method = &m
				return
			}
		}

	})
	return init.method

}
