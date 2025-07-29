// The life cycle of a service is controlled by the container.
package vdi

type Transient[TOwner any, T any] struct {
	Owner interface{}
	Init  func(owner *TOwner) T
	fn    func(owner *TOwner, args ...any) T
}

func (s *Transient[TOwner, T]) Get(args ...any) T {
	if s.Owner == nil {
		panic("Transient[TOwner, T] requires an owner")
	}
	if s.fn != nil {
		return s.fn(s.Owner.(*TOwner), args...)
	}
	return s.Init(s.Owner.(*TOwner))
}

func (s *Transient[TOwner, T]) SetFunc(fn func(owner *TOwner, args ...any) T) {
	s.fn = fn
}
