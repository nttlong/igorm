package services

import (
	"net/http"
	"reflect"
)

var NewServiceContext func(req *http.Request, res http.ResponseWriter) interface{}

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
