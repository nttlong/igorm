package vdi

import (
	"fmt"
	"reflect"
	"sync"
)

type Singleton[TOwner any, T any] struct {
	value T
	Owner interface{}
	Init  func(owner *TOwner) T
	once  sync.Once
}

func (s *Singleton[TOwner, T]) Set(value T) {
	s.value = value

}
func (s *Singleton[TOwner, T]) Get() T {
	if s.Owner == nil {
		panic("Singleton[TOwner, T] requires an owner")
	}
	s.once.Do(func() {
		typ := reflect.TypeOf(s.Owner)
		if typ.Kind() == reflect.Ptr {
			s.value = s.Init(s.Owner.(*TOwner))
		} else {
			owner := s.Owner.(TOwner)
			if s.Init == nil {
				ownerType := reflect.TypeFor[TOwner]()
				valType := reflect.TypeFor[T]()
				panic(fmt.Errorf("Singleton[%s, %s] requires an Init function", ownerType.String(), valType.String()))
			}
			s.value = s.Init(&owner)
		}

	})
	return s.value
}
