package vapi

import (
	"encoding/json"
	"net/http"
	"reflect"
)

func (web *webHandlerRunnerType) ExecGet(handler webHandler, w http.ResponseWriter, r *http.Request) error {
	args := make([]reflect.Value, handler.apiInfo.Method.Type.NumIn())
	var err error
	args[0], err = web.ResolveReceiverValue(handler)
	if err != nil {
		return err
	}

	args[handler.apiInfo.IndexOfArg], err = web.CreateHttpContext(handler, w, r)
	if err != nil {
		return err
	}

	retArgs := web.MethodCall(handler, args)
	if len(retArgs) > 0 {
		if err, ok := retArgs[len(retArgs)-1].Interface().(error); ok {
			return err
		}
		if len(retArgs) > 1 {
			retIntefaces := []interface{}{}
			for i := 0; i < len(retArgs)-1; i++ {
				retIntefaces = append(retIntefaces, retArgs[i].Interface())
			}

			retArgs = retArgs[0 : len(retArgs)-2]
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if len(retIntefaces) == 1 {
				json.NewEncoder(w).Encode(retIntefaces[0])

			} else {
				json.NewEncoder(w).Encode(retIntefaces)
			}
			return nil

		}
	}

	return nil
}
