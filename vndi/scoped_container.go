package vdi

import (
	"fmt"
	"reflect"
)

type ScopedContainer struct {
	parent    *RootContainer
	instances map[reflect.Type]any
}

func (c *ScopedContainer) ResolveByType(t reflect.Type) (any, error) {
	if c.parent.IsTransient(t) {
		factory, ok := c.parent.GetFactory(t)
		if !ok {
			return nil, fmt.Errorf("type %v not registered", t)
		}
		return callFactoryMust(factory, c), nil
	}

	if val, ok := c.instances[t]; ok {
		return val, nil
	}

	factory, ok := c.parent.GetFactory(t)
	if !ok {
		return nil, fmt.Errorf("type %v not registered", t)
	}

	result := callFactoryMust(factory, c)
	c.instances[t] = result
	return result, nil
}
