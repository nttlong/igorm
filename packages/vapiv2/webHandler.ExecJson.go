package vapi

import (
	"encoding/json"
	"net/http"
	"reflect"
)

func (web *webHandlerRunnerType) ExecJson(handler webHandler, w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return err
	}
	ReceiverValue, err := web.ResolveReceiverValue(handler)
	if err != nil {
		return err
	}

	var bodyData reflect.Value

	// Duyệt tất cả key/value trong form
	if handler.apiInfo.IndexOfRequestBody > -1 {
		bodyData = reflect.New(handler.apiInfo.TypeOfRequestBodyElem)

		if err := json.NewDecoder(r.Body).Decode(bodyData.Interface()); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return err
		}

	}
	args := make([]reflect.Value, handler.apiInfo.Method.Type.NumIn())
	args[0] = ReceiverValue
	if handler.apiInfo.IndexOfRequestBody > -1 {
		args[handler.apiInfo.IndexOfRequestBody] = bodyData

	}

	if handler.apiInfo.ReceiverTypeElem.Kind() == reflect.Ptr {
		handler.apiInfo.ReceiverTypeElem = handler.apiInfo.ReceiverTypeElem.Elem()
	}

	context := reflect.New(handler.apiInfo.TypeOfArgsElem)

	context.Elem().FieldByName("Req").Set(reflect.ValueOf(r))
	context.Elem().FieldByName("Res").Set(reflect.ValueOf(w))

	args[handler.apiInfo.IndexOfArg] = context
	for i := 1; i < handler.apiInfo.Method.Type.NumIn(); i++ {
		if handler.apiInfo.Method.Type.In(i).Kind() != reflect.Ptr {
			if args[i].Kind() == reflect.Ptr {
				args[i] = args[i].Elem()

			}
		}

	}

	retArgs := handler.apiInfo.Method.Func.Call(args)
	if len(retArgs) > 0 {
		if err, ok := retArgs[len(retArgs)-1].Interface().(error); ok {
			return err
		}
		if len(retArgs) > 2 {
			retIntefaces := []interface{}{}
			for i := 0; i < len(retArgs)-1; i++ {
				retIntefaces = append(retIntefaces, retArgs[i].Interface())
			}

			retArgs = retArgs[0 : len(retArgs)-2]
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(retIntefaces)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(retArgs[0].Elem().Interface())
		}
		// Ví dụ: trả về dạng JSON

	}

	return nil
}
