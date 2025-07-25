package vdi

import (
	"fmt"
	"reflect"
)

func InjectFields(container Container, target any) error {
	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return fmt.Errorf("target must be a non-nil pointer")
	}
	val = val.Elem()
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("inject")
		if tag == "" {
			continue
		}
		if !val.Field(i).CanSet() {
			return fmt.Errorf("cannot inject into unexported field: %s", field.Name)
		}

		// Resolve dependency
		depType := field.Type
		resolved, err := container.ResolveByType(depType)
		if err != nil {
			return fmt.Errorf("injecting %s: %w", field.Name, err)
		}
		val.Field(i).Set(reflect.ValueOf(resolved))
	}
	return nil
}
