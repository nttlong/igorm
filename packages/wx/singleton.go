package wx

import (
	"fmt"
	"reflect"
	"sync"
)

type Singleton[T any] struct {
	// Owner *TService
	once sync.Once
	ins  *T
	err  error

	init func() (*T, error)
}

func (s *Singleton[T]) GetInstance() *T {

	s.once.Do(func() {
		if s.init == nil {
			panic(fmt.Sprintf("%s not initialized,please call Init() of %s first", reflect.TypeOf(s).String(), reflect.TypeOf(s).String()))
		}
		r, err := s.init()
		if err != nil {
			panic(err)
		}
		s.ins = r
	})
	return s.ins
}
func (s *Singleton[T]) Init(fn func() (*T, error)) {
	s.init = fn
}
