package vapi

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"unsafe"
	vapiErr "vapi/errors"
)

func (svc *serviceUtilsType) IsFieldSingleton(field reflect.StructField) bool {
	typ := field.Type
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.PkgPath() != svc.pkgPath {
		return false

	}

	return strings.HasPrefix(typ.String(), svc.checkSingletonTypeName)
}
func (svc *serviceUtilsType) IsFieldScoped(field reflect.StructField) bool {
	typ := field.Type
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.PkgPath() != svc.pkgPath {
		return false

	}

	return strings.HasPrefix(typ.String(), svc.checkScopeTypeName)

}
func (svc *serviceUtilsType) IsSingletonType(typ reflect.Type) bool {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.PkgPath() != svc.pkgPath {
		return false
	}

	return strings.HasPrefix(typ.String(), svc.checkSingletonTypeName)
}

type initNewSingletonByType struct {
	once     sync.Once
	instance reflect.Value
}

var initNewSingletonByTypeCache = sync.Map{}

func (svc *serviceUtilsType) NewSingletonByType(typ reflect.Type) reflect.Value {
	actual, _ := initNewSingletonByTypeCache.LoadOrStore(typ, &initNewSingletonByType{})
	initService := actual.(*initNewSingletonByType)
	initService.once.Do(func() {
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		initService.instance = reflect.New(typ)
	})
	return initService.instance
}
func (svc *serviceUtilsType) CreateTransient(receiverValue *reflect.Value, field reflect.StructField) {
	fieldValue := receiverValue.Elem().FieldByIndex(field.Index)
	if fieldValue.Kind() == reflect.Ptr {
		fieldValue = fieldValue.Elem()
	}
	fieldType := field.Type
	if fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
	}

	instanceOfField := reflect.New(fieldType)

	if fieldValue.Kind() == reflect.Ptr {
		fieldValue = fieldValue.Elem()
	}

	fieldValue.Set(instanceOfField.Elem())
}

func (svc *serviceUtilsType) CreateSingeton(receiverValue *reflect.Value, field reflect.StructField) {
	fieldValue := receiverValue.Elem().FieldByIndex(field.Index)
	if fieldValue.Kind() == reflect.Ptr {
		fieldValue = fieldValue.Elem()
	}
	fieldType := field.Type
	if fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
	}

	instanceOfField := svc.NewSingletonByType(fieldType)
	//instanceOfField := reflect.New(fieldType)

	if fieldValue.Kind() == reflect.Ptr {
		fieldValue = fieldValue.Elem()
	}

	fieldValue.Set(instanceOfField.Elem())
}
func (svc *serviceUtilsType) IsInjector(typ reflect.Type) bool {
	return svc.isInjectorInternal(typ, make(map[reflect.Type]struct{}))
}

func (svc *serviceUtilsType) isInjectorInternal(typ reflect.Type, visited map[reflect.Type]struct{}) bool {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return false
	}

	// Nếu đã kiểm tra rồi thì bỏ qua để tránh vòng lặp
	if _, ok := visited[typ]; ok {
		return false
	}
	visited[typ] = struct{}{}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldType := field.Type
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}
		if fieldType.Kind() != reflect.Struct {
			continue
		}
		if svc.IsFieldSingleton(field) || svc.IsFieldScoped(field) {
			return true
		}
		if svc.isInjectorInternal(fieldType, visited) {
			return true
		}
	}
	return false
}

type initNewService struct {
	once     sync.Once
	instance *reflect.Value
	err      error
}

var initNewServiceCache = sync.Map{}

func (svc *serviceUtilsType) NewServiceOneTime(typ reflect.Type) (*reflect.Value, error) {
	actual, _ := initNewServiceCache.LoadOrStore(typ, &initNewService{})
	initService := actual.(*initNewService)
	initService.once.Do(func() {
		initService.instance, initService.err = svc.NewService(typ, nil, nil)
	})
	if initService.err != nil {
		return nil, initService.err
	}

	return initService.instance, nil

}

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
type initGetServiceInfo struct {
	once     sync.Once
	instance *serviceRecord
	err      error
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

var initGetServiceInfoCache = sync.Map{}

func (svc *serviceUtilsType) GetServiceInfo(typ reflect.Type) (*serviceRecord, error) {
	actual, _ := initGetServiceInfoCache.LoadOrStore(typ, &initGetServiceInfo{})
	initService := actual.(*initGetServiceInfo)
	initService.once.Do(func() {
		initService.instance, initService.err = svc.getServiceInfo(typ)
	})
	return initService.instance, initService.err

}

func (svc *serviceUtilsType) NewService(typ reflect.Type, req *http.Request, res http.ResponseWriter) (*reflect.Value, error) {
	info, err := svc.GetServiceInfo(typ)
	if err != nil {
		return nil, err
	}

	ret := reflect.New(info.ReciverTypeElem) //<-- info.ReciverType is always ptr
	retEle := ret.Elem()
	for i, val := range info.SingletonValue {
		field := retEle.FieldByIndex(info.SingletonFieldIndex[i])
		if field.Kind() == reflect.Ptr {
			if field.IsNil() {
				field.Set(val)
				continue
			}

		}
		field.Set(val.Elem())

	}

	for _, fieldIndex := range info.ScopedFieldIndex { //<-- hay sua lai bang cach dunh unsafe pionter
		field := retEle.FieldByIndex(fieldIndex)

		val := svc.CreateScope(field.Type())
		scvContext := NewServiceContext(req, res)
		val.Elem().FieldByName("Ctx").Set(reflect.ValueOf(scvContext))

		if field.Kind() == reflect.Ptr {
			if field.IsNil() {
				field.Set(val)
				continue
			}

		}
		field.Set(val.Elem())
	}
	retVal := info.NewMethod.Func.Call([]reflect.Value{ret})
	if retVal[0].Interface() != nil {
		return nil, retVal[0].Interface().(error)
	}
	return &ret, nil

}
func (svc *serviceUtilsType) LoadSingletonFields(serviceVal *reflect.Value) error {
	for i := 0; i < serviceVal.Elem().NumField(); i++ {
		field := serviceVal.Elem().Type().Field(i)
		if svc.IsFieldSingleton(field) {
			svc.CreateSingeton(serviceVal, field)

		}

	}
	return nil

}
func (svc *serviceUtilsType) LoadScopedFields(serviceVal *reflect.Value) error {
	for i := 0; i < serviceVal.Elem().NumField(); i++ {
		field := serviceVal.Elem().Type().Field(i)
		if svc.IsFieldScoped(field) {
			svc.CreateTransient(serviceVal, field)

		}

	}
	return nil

}

func (svc *serviceUtilsType) setFieldUnsafe(field, value reflect.Value) {
	// Tạo một reflect.Value trỏ trực tiếp vào ô nhớ field
	fv := reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()

	if field.Kind() == reflect.Ptr {
		fv.Set(value)
	} else {
		fv.Set(value.Elem())
	}
}

func (svc *serviceUtilsType) NewServiceOptimize(typ reflect.Type) (*reflect.Value, error) {
	info, err := svc.GetServiceInfo(typ)
	if err != nil {
		return nil, err
	}

	// Tạo instance mới
	ret := reflect.New(info.ReciverType.Elem())
	basePtr := unsafe.Pointer(ret.Pointer()) // trỏ tới vùng nhớ của struct

	// Inject Singleton bằng offset
	for i, val := range info.SingletonValue {
		fieldPtr := unsafe.Pointer(uintptr(basePtr) + info.SingletonOffsets[i])
		fieldVal := reflect.NewAt(info.SingletonTypes[i], fieldPtr).Elem()
		fieldVal.Set(val) //<--panic: reflect.Set: value of type *vapi.Singleton[vapi.FileUtils] is not assignable to type vapi.Singleton[vapi.FileUtils]

	}

	// Inject Scoped bằng offset
	for i, fieldType := range info.ScopedTypes {
		newScoped := svc.CreateScope(fieldType)
		fieldPtr := unsafe.Pointer(uintptr(basePtr) + info.ScopedOffsets[i])
		fieldVal := reflect.NewAt(fieldType, fieldPtr).Elem()
		fieldVal.Set(newScoped)
	}

	// Gọi hàm New()
	retVal := info.NewMethod.Func.Call([]reflect.Value{ret})
	if retVal[0].Interface() != nil {
		return nil, retVal[0].Interface().(error)
	}

	return &ret, nil
}
