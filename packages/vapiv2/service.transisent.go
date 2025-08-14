package vapi

import (
	"fmt"
	"reflect"
)

type Transient[T any, TService any] struct {
	Owner *TService
	ins   *T
	err   error
	init  func() (*T, error)
}

func (t *Transient[T, TService]) Init(fn func() (*T, error)) {
	t.init = fn
}
func (t *Transient[T, TService]) GetInstance() (*T, error) {
	if t.init == nil {
		return nil, fmt.Errorf("%s not initialized,please call Init() of %s first", reflect.TypeOf(t).String(), reflect.TypeOf(t).String())
	}
	r, err := t.init()
	return r, err

}
