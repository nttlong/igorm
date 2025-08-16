package services

import (
	"fmt"
	"reflect"
	"sync"
	vapiErr "wx/errors"
)

type serviceRecord struct {
	SingletonFieldIndex [][]int
	ScopedFieldIndex    [][]int

	ReciverType     reflect.Type
	ReciverTypeElem reflect.Type

	SingletonOffsets []uintptr
	SingletonTypes   []reflect.Type
	SingletonValue   []reflect.Value

	ScopedOffsets []uintptr
	ScopedTypes   []reflect.Type

	NewMethod reflect.Method
}

func (svc *serviceUtilsType) getServiceInfo(typ reflect.Type) (*serviceRecord, error) {
	var newMethod reflect.Method
	foungNewMethod := false
	ptrType := typ
	if ptrType.Kind() != reflect.Ptr {
		ptrType = reflect.PointerTo(typ)
	}

	for i := 0; i < ptrType.NumMethod(); i++ {
		if ptrType.Method(i).Name == "New" {
			newMethod = ptrType.Method(i)
			foungNewMethod = true

			break
		}
	}
	if foungNewMethod {
		SingletonOffsets, SingletonTypes := svc.GetSingletonFieldsOffsetPtr(typ)
		ScopedOffsets, ScopedTypes := svc.GetGetScopeFieldsOffsetPtr(typ)
		typeEle := typ
		if typeEle.Kind() == reflect.Ptr {
			typeEle = typeEle.Elem()
		}

		ret := &serviceRecord{
			NewMethod:           newMethod,
			ReciverType:         typ,
			ReciverTypeElem:     typeEle,
			SingletonFieldIndex: svc.getSingletonFieldsInternal(typ, make(map[reflect.Type]bool)),
			ScopedFieldIndex:    svc.getScopedFieldsInternal(typ, make(map[reflect.Type]bool)),
			SingletonOffsets:    SingletonOffsets,
			SingletonTypes:      SingletonTypes,
			ScopedOffsets:       ScopedOffsets,
			ScopedTypes:         ScopedTypes,
			SingletonValue:      []reflect.Value{},
		}
		ft := typ
		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}
		for _, x := range ret.SingletonFieldIndex { // vi la singleton nen tao luon cho nay
			val := svc.CreateSingletonInstance(ft.FieldByIndex(x).Type)
			ret.SingletonValue = append(ret.SingletonValue, val)

		}
		return ret, nil

	} else {
		errMsg := fmt.Sprintf("New function was not found in %s. injector need New function", typ.String())
		return nil, vapiErr.NewServiceInitError(errMsg)
	}

}

type initGetServiceInfo struct {
	once     sync.Once
	instance *serviceRecord
	err      error
}

var initGetServiceInfoCache = sync.Map{}

func (svc *serviceUtilsType) GetServiceInfo(typ reflect.Type) (*serviceRecord, error) {
	actual, _ := initGetServiceInfoCache.LoadOrStore(typ, &initGetServiceInfo{})
	initService := actual.(*initGetServiceInfo)
	initService.once.Do(func() {
		initService.instance, initService.err = svc.getServiceInfo(typ)
	})
	return initService.instance, initService.err

}
