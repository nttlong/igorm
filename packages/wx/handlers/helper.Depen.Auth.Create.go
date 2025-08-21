package handlers

import (
	"net/http"
	"reflect"
)

func (h *helperType) DepenAuthCreateInstance(newMethod reflect.Method, ctxVale reflect.Value) (*reflect.Value, error) {
	r := reflect.New(newMethod.Type.In(0).Elem())
	arsg := make([]reflect.Value, newMethod.Type.NumIn())
	arsg[0] = r
	arsg[1] = ctxVale
	res := newMethod.Func.Call(arsg)
	if res[0].Interface() != nil {
		return nil, res[0].Interface().(error)
	}
	return &r, nil

}

var DepenAuthCreateCreateAuthContext func(req *http.Request, Res http.ResponseWriter) (*reflect.Value, error)

func (h *helperType) DepenAuthCreate(typ reflect.Type, req *http.Request, Res http.ResponseWriter) (*reflect.Value, error) {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()

	}
	r := reflect.New(typ)
	contextvalue, err := DepenAuthCreateCreateAuthContext(req, Res)
	if err != nil {
		return nil, err
	}
	r.Elem().FieldByName("Context").Set(*contextvalue)

	//fmt.Println(typ.String())

	return &r, nil
}
