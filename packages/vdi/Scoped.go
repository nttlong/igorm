package vdi

// The life cycle of a service is controlled by the owner.
type Scoped[TOwner any, T any] struct {
	Value T
	Owner interface{}
	Init  func(owner *TOwner) T
}

func (s *Scoped[TOwner, T]) Get() T {
	if s.Owner == nil {
		panic("Scoped[TOwner, T] requires an owner")
	}
	if s.Init == nil {
		return s.Value
	}
	s.Value = s.Init(s.Owner.(*TOwner))
	s.Init = nil
	return s.Value
}
