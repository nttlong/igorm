package wx

import (
	"net/http"
	"reflect"
	"wx/handlers"
)

type Handler struct {
	Req        *http.Request
	Res        http.ResponseWriter
	schema     string
	rootAbsUrl string
}

func handlerIsArgHandler(typ reflect.Type, visited map[reflect.Type]bool) ([]int, []int, []int) {
	if typ == nil || visited[typ] {
		return nil, nil, nil
	}
	visited[typ] = true
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ == reflect.TypeOf(Handler{}) || typ.ConvertibleTo(reflect.TypeOf(Handler{})) {
		var fieldReq reflect.StructField
		var fieldRes reflect.StructField
		var ok bool
		if fieldReq, ok = typ.FieldByName("Req"); !ok {
			return nil, nil, nil
		}
		if fieldRes, ok = typ.FieldByName("Res"); !ok {
			return nil, nil, nil

		}
		return []int{}, fieldReq.Index, fieldRes.Index
	}

	if typ.Kind() != reflect.Struct {
		return nil, nil, nil
	}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if indexes, _, _ := handlerIsArgHandler(field.Type, visited); indexes != nil {
			fieldIndex := append(field.Index, indexes...)
			var fieldReq reflect.StructField
			var fieldRes reflect.StructField
			var ok bool
			if fieldReq, ok = typ.FieldByName("Req"); !ok {
				return nil, nil, nil
			}
			if fieldRes, ok = typ.FieldByName("Res"); !ok {
				return nil, nil, nil

			}
			return fieldIndex, fieldReq.Index, fieldRes.Index
		}
	}
	return nil, nil, nil
}

func createHandlerContext(typ reflect.Type, reqIndex []int, resIndex []int, r *http.Request, w http.ResponseWriter) reflect.Value {
	if typ == reflect.TypeOf(Handler{}) {
		return reflect.ValueOf(&Handler{
			Req: r,
			Res: w,
		})
	}

	ret := reflect.New(typ)
	ret.Elem().FieldByIndex(reqIndex).Set(reflect.ValueOf(r))
	ret.Elem().FieldByIndex(resIndex).Set(reflect.ValueOf(w))
	return ret
}
func init() {
	handlers.HandlerIsArgHandler = handlerIsArgHandler
	handlers.CreateHandlerContext = createHandlerContext
}
