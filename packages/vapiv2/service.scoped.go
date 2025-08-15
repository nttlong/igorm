package vapi

// import (
// 	"fmt"
// 	"reflect"
// 	"sync"
// )

// type Scoped[T any] struct {
// 	ins  *T
// 	err  error
// 	once sync.Once
// 	init func() (*T, error)
// }

// func (s *Scoped[T]) GetInstance() *T {

// 	s.once.Do(func() {
// 		if s.init == nil {
// 			panic(fmt.Sprintf("%s not initialized,please call Init() of %s first", reflect.TypeOf(s).String(), reflect.TypeOf(s).String()))
// 		}
// 		r, err := s.init()
// 		if err != nil {
// 			panic(err)
// 		}
// 		s.ins = r
// 	})
// 	return s.ins
// }
// func (s *Scoped[T]) Init(fn func() (*T, error)) {
// 	s.init = fn
// }
