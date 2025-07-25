package vdi

import (
	"fmt"
	"reflect"
	"sync"
)

// NamedKey là tổ hợp giữa reflect.Type và tên tuỳ chọn
type NamedKey struct {
	Type reflect.Type
	Name string
}

// NamedRegistry lưu trữ các service đã đăng ký kèm tên
type NamedRegistry struct {
	mu    sync.RWMutex
	store map[NamedKey]any
}

func NewNamedRegistry() *NamedRegistry {
	return &NamedRegistry{
		store: make(map[NamedKey]any),
	}
}

func (r *NamedRegistry) Register(t reflect.Type, name string, instance any) error {
	key := NamedKey{Type: t, Name: name}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.store[key]; exists {
		return fmt.Errorf("service of type %s with name '%s' already registered", t.String(), name)
	}
	r.store[key] = instance
	return nil
}

func (r *NamedRegistry) Resolve(t reflect.Type, name string) (any, error) {
	key := NamedKey{Type: t, Name: name}

	r.mu.RLock()
	defer r.mu.RUnlock()

	if inst, ok := r.store[key]; ok {
		return inst, nil
	}
	return nil, fmt.Errorf("no service found for type %s with name '%s'", t.String(), name)
}
