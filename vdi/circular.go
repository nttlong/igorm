package vdi

import (
	"fmt"
	"reflect"
)

type resolveStack struct {
	stack []reflect.Type
}

func (r *resolveStack) Push(t reflect.Type) error {
	for _, s := range r.stack {
		if s == t {
			return fmt.Errorf("circular dependency detected: %s", r.StringWith(t))
		}
	}
	r.stack = append(r.stack, t)
	return nil
}

func (r *resolveStack) Pop() {
	if len(r.stack) > 0 {
		r.stack = r.stack[:len(r.stack)-1]
	}
}

func (r *resolveStack) StringWith(t reflect.Type) string {
	chain := ""
	for _, s := range r.stack {
		chain += s.String() + " -> "
	}
	return chain + t.String()
}
