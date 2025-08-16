package services

import (
	"reflect"
	"sync"
)

type initGetScopedOffsets struct {
	once    sync.Once
	offsets []uintptr
	types   []reflect.Type
}

var cacheGetScopedOffsets = sync.Map{}

// Trả về offsets và types của các field Scoped trong struct
func (svc *serviceUtilsType) GetGetScopeFieldsOffsetPtr(typ reflect.Type) ([]uintptr, []reflect.Type) {
	actual, _ := cacheGetScopedOffsets.LoadOrStore(typ, &initGetScopedOffsets{})
	initService := actual.(*initGetScopedOffsets)
	initService.once.Do(func() {
		initService.offsets, initService.types = svc.getScopedOffsetsInternal(typ, map[reflect.Type]bool{})
	})
	return initService.offsets, initService.types
}

func (svc *serviceUtilsType) getScopedOffsetsInternal(typ reflect.Type, visited map[reflect.Type]bool) ([]uintptr, []reflect.Type) {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil, nil
	}
	if visited[typ] {
		return nil, nil
	}
	visited[typ] = true

	var offsets []uintptr
	var types []reflect.Type

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		if svc.IsFieldScoped(field) {
			offsets = append(offsets, field.Offset)
			types = append(types, field.Type)
		} else {
			// Đệ quy xuống struct con
			subOffsets, subTypes := svc.getScopedOffsetsInternal(field.Type, visited)
			for j := range subOffsets {
				offsets = append(offsets, field.Offset+subOffsets[j])
				types = append(types, subTypes[j])
			}
		}
	}
	return offsets, types
}
