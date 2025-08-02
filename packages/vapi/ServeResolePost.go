package vapi

import (
	"encoding/json"
	"net/http"
	"reflect"
)

func serverResolverPost(w http.ResponseWriter, r *http.Request, method reflect.Method, path string) {
	inputParam := mapInputParamInfo[path]
	// inputValues := make([]reflect.Value, method.Type.NumIn())

	receiver := mapInstanceInit[path].Call([]reflect.Value{})[0]
	applyContext(path, receiver, w, r)
	decoder := json.NewDecoder(r.Body)
	inputType := method.Type.In(1)
	if inputType.Kind() == reflect.Ptr {
		inputType = inputType.Elem()
	}
	data := reflect.New(inputType).Interface()
	err := decoder.Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	inputPutValues := make([]reflect.Value, method.Type.NumIn())
	inputPutValues[0] = receiver
	if inputParam.dataParamIndex != -1 {
		inputPutValues[inputParam.dataParamIndex] = reflect.ValueOf(data).Elem()

	}
	if inputParam.userParamIndex != -1 {
		inputPutValues[inputParam.userParamIndex] = reflect.ValueOf(UserClaims{})
	}
	output := method.Func.Call(inputPutValues)
	w.Header().Set("Content-Type", "application/json")

	outPutValues := make([]interface{}, len(output))

	for i, v := range output {
		outPutValues[i] = v.Interface()
	}

	jsonData, _ := json.Marshal(outPutValues[0])

	w.Write(jsonData)
}
