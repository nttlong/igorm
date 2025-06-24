package unvscore

// import (
// 	"context"
// 	"dbx"
// 	"reflect"
// 	"strings"
// 	"sync"

// 	cacher "unvs.core/cacher"
// 	config "unvs.core/config"
// )

// type BaseService[TContainer any] struct {
// 	Container TContainer
// }
// type Singleton[TOwner any, T any] struct {
// 	Value T
// 	Owner interface{}
// 	Init  func(owner TOwner) T
// 	once  sync.Once
// }

// // The life cycle of a service is controlled by the owner.
// type Scoped[TOwner any, T any] struct {
// 	Value T
// 	Owner interface{}
// 	Init  func(owner TOwner) T
// }

// // The life cycle of a service is controlled by the container.
// type Transient[TOwner any, T any] struct {
// 	Value T
// 	Owner interface{}
// 	OnGet func(owner TOwner) T
// }

// func (s *Transient[TOwner, T]) Get() T {
// 	if s.Owner == nil {
// 		panic("Transient[TOwner, T] requires an owner")
// 	}
// 	return s.OnGet(*s.Owner.(*TOwner))
// }

// func (s *Singleton[TOwner, T]) Get() T {
// 	if s.Owner == nil {
// 		panic("Singleton[TOwner, T] requires an owner")
// 	}
// 	s.once.Do(func() {
// 		s.Value = s.Init(*s.Owner.(*TOwner))
// 	})
// 	return s.Value
// }

// func (s *Scoped[TOwner, T]) Get() T {
// 	if s.Owner == nil {
// 		panic("Scoped[TOwner, T] requires an owner")
// 	}
// 	if s.Init == nil {
// 		return s.Value
// 	}
// 	s.Value = s.Init(*s.Owner.(*TOwner))
// 	s.Init = nil
// 	return s.Value
// }

// type BaseServiceInfo struct {
// 	Tenant      string
// 	Lang        string
// 	AccessToken string
// 	FeatureId   string
// 	Context     context.Context
// }
// type sampleService struct{}

// func Resolve[T any](params any) *T {
// 	typ := reflect.TypeFor[T]()
// 	if typ.Kind() == reflect.Ptr {
// 		typ = typ.Elem()
// 	}
// 	if typ.Kind() != reflect.Struct {
// 		panic("CreateService[T] requires a struct type")
// 	}
// 	retVal := reflect.New(typ).Elem()
// 	packagePath := reflect.TypeOf(sampleService{}).PkgPath()

// 	for i := 0; i < typ.NumField(); i++ {
// 		field := typ.Field(i)
// 		// get field value

// 		if field.Type.Kind() == reflect.Struct {
// 			fieldType := field.Type

// 			fieldTypeName := fieldType.Name()
// 			fieldPackagePath := fieldType.PkgPath()
// 			if fieldPackagePath != packagePath {
// 				continue
// 			}
// 			isInjector := false
// 			isInjector = isInjector || strings.HasPrefix(fieldTypeName, "Singleton[")
// 			isInjector = isInjector || strings.HasPrefix(fieldTypeName, "Scoped[")
// 			isInjector = isInjector || strings.HasPrefix(fieldTypeName, "Transient[")

// 			if isInjector {
// 				val := retVal.Field(i)
// 				ownerField := val.FieldByName("Owner")
// 				if ownerField.IsValid() && ownerField.CanSet() {
// 					ownerField.Set(retVal.Addr())
// 				}

// 			}

// 		}

// 	}
// 	retInterface := retVal.Addr().Interface()
// 	return retInterface.(*T)
// }
