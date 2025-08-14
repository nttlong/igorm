package vapi

import (
	"fmt"
	"reflect"
	"sync"
)

type serviceInfo struct {
	reciverType reflect.Type
	fn          reflect.Value
}

var cacheRegisterService = map[reflect.Type]serviceInfo{}

func RegisterService[TService any](fn func(service *TService) error) {
	cacheRegisterService[reflect.TypeFor[TService]()] = serviceInfo{
		fn:          reflect.ValueOf(fn),
		reciverType: reflect.TypeFor[TService](),
	}
}

type initCreatService struct {
	once     sync.Once
	instance interface{}
	err      error
}

var initRegisterServiceCache = sync.Map{}

func Service[T any]() (*T, error) {
	actual, _ := initRegisterServiceCache.LoadOrStore(reflect.TypeFor[T](), &initCreatService{})
	initService := actual.(*initCreatService)
	initService.once.Do(func() {
		initService.instance, initService.err = creatService[T]()
	})
	if initService.err != nil {
		return nil, initService.err
	}
	return initService.instance.(*T), nil

}
func creatService[T any]() (*T, error) {
	typ := reflect.TypeFor[T]()
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if svInfo, ok := cacheRegisterService[typ]; ok {
		receiverValue := reflect.New(svInfo.reciverType)
		for i := 0; i < svInfo.reciverType.NumField(); i++ {
			field := svInfo.reciverType.Field(i)

			if serviceUtils.IsFieldSingleton(field) {
				serviceUtils.CreateSingeton(&receiverValue, field)
			} else if serviceUtils.IsFieldTransient(field) {
				serviceUtils.CreateSingeton(&receiverValue, field)
				fmt.Println(field.Name)
			}

		}
		ret := svInfo.fn.Call([]reflect.Value{receiverValue})
		if ret[0].Interface() != nil {
			return nil, ret[0].Interface().(error)
		} else {
			return receiverValue.Interface().(*T), nil
		}

	}
	return nil, fmt.Errorf("service %s not found", typ.String())

}
