package services

import (
	"reflect"
	"sync"
)

type initGetSingletonOffsets struct {
	once    sync.Once
	offsets []uintptr
	types   []reflect.Type
}

var cacheGetSingletonOffsets = sync.Map{}

// Trả về offsets và types của các field Singleton trong struct
func (svc *serviceUtilsType) GetSingletonFieldsOffsetPtr(typ reflect.Type) ([]uintptr, []reflect.Type) {
	actual, _ := cacheGetSingletonOffsets.LoadOrStore(typ, &initGetSingletonOffsets{})
	initService := actual.(*initGetSingletonOffsets)
	initService.once.Do(func() {
		initService.offsets, initService.types = svc.getSingletonOffsetsInternal(typ, map[reflect.Type]bool{})
	})
	return initService.offsets, initService.types
}

func (svc *serviceUtilsType) getSingletonOffsetsInternal(typ reflect.Type, visited map[reflect.Type]bool) ([]uintptr, []reflect.Type) {
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

		if svc.IsFieldSingleton(field) {
			offsets = append(offsets, field.Offset)
			types = append(types, field.Type)
		} else {
			// Đệ quy nếu là struct
			subOffsets, subTypes := svc.getSingletonOffsetsInternal(field.Type, visited)
			for j := range subOffsets {
				offsets = append(offsets, field.Offset+subOffsets[j])
				types = append(types, subTypes[j])
			}
		}
	}
	return offsets, types
}
