// The life cycle of a service is controlled by the container.
package vdi

type Transient[TOwner any, T any] struct {
	Value *T
	Owner interface{}
	Init  func(owner *TOwner) T
}

func (s *Transient[TOwner, T]) Get() T {
	if s.Owner == nil {
		panic("Transient[TOwner, T] requires an owner")
	}
	return s.Init(s.Owner.(*TOwner))
}
