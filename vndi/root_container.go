package vdi

import (
	"fmt"
	"reflect"
	"sync"
)

type RootContainer struct {
	typeFactories  map[reflect.Type]func(Container) any
	transientTypes map[reflect.Type]struct{}
	mu             sync.RWMutex
}

func NewRootContainer() *RootContainer {
	return &RootContainer{
		typeFactories:  make(map[reflect.Type]func(Container) any),
		transientTypes: make(map[reflect.Type]struct{}),
	}
}

func (c *RootContainer) RegisterSingleton(factory any) {
	t := factoryType(factory)
	var once sync.Once
	var instance any
	c.mu.Lock()
	defer c.mu.Unlock()

	c.typeFactories[t] = func(_ Container) any {
		once.Do(func() {
			instance = callFactoryMust(factory, nil)
		})
		return instance
	}
}

func (c *RootContainer) RegisterScoped(factory any) {
	t := factoryType(factory)
	c.mu.Lock()
	defer c.mu.Unlock()

	c.typeFactories[t] = func(scope Container) any {
		return callFactoryMust(factory, scope)
	}
}

func (c *RootContainer) RegisterTransient(factory any) {
	t := factoryType(factory)
	c.mu.Lock()
	defer c.mu.Unlock()

	c.typeFactories[t] = func(scope Container) any {
		return callFactoryMust(factory, scope)
	}
	c.transientTypes[t] = struct{}{}
}

func (c *RootContainer) CreateScope() Scope {
	return &ScopedContainer{
		parent:    c,
		instances: make(map[reflect.Type]any),
	}
}

func (c *RootContainer) ResolveByType(t reflect.Type) (any, error) {
	factory, ok := c.GetFactory(t)
	if !ok {
		return nil, fmt.Errorf("type %v not registered", t)
	}
	return factory(c), nil
}

func (c *RootContainer) GetFactory(t reflect.Type) (func(Container) any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	f, ok := c.typeFactories[t]
	return f, ok
}

func (c *RootContainer) IsTransient(t reflect.Type) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.transientTypes[t]
	return ok
}
