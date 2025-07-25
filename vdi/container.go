package vdi

import "reflect"

type Container interface {
	ResolveByType(t reflect.Type) (any, error)
}

type Scope interface {
	Container
}

type Root interface {
	Container
	RegisterSingleton(factory any)
	RegisterScoped(factory any)
	RegisterTransient(factory any)
	CreateScope() Scope
}
