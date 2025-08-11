package vapi

import (
	"net/http"
	"reflect"
)

func (web *webHandlerRunnerType) ResolveReceiverValue(handler webHandler, r *http.Request) (reflect.Value, error) {
	key := handler.apiInfo.ReceiverType.String() + "/webHandlerRunnerType/ResolveReceiverValue"
	ret, err := OnceCall(key, func() (*reflect.Value, error) {
		result := handler.initFunc.Call([]reflect.Value{})
		if result[1].IsValid() && !result[1].IsNil() {
			return nil, result[1].Interface().(error)
		}
		instanceType := handler.apiInfo.ReceiverType
		if instanceType.Kind() == reflect.Ptr {
			instanceType = instanceType.Elem()
		}
		baseUrlField, ok := instanceType.FieldByName("BaseUrl")
		if ok {
			fieldIndex := baseUrlField.Index
			if len(fieldIndex) > 1 {
				parentFieldIndex := fieldIndex[0 : len(fieldIndex)-1]
				parentBaseUrlField := instanceType.FieldByIndex(parentFieldIndex)
				parentBaseurlFieldType := parentBaseUrlField.Type
				if parentBaseurlFieldType.Kind() == reflect.Ptr {
					parentBaseurlFieldType = parentBaseurlFieldType.Elem()
				}
				if parentBaseurlFieldType == reflect.TypeOf(Service{}) {
					instanceValue := &result[0]
					_, _, _, baseURL := getBaseURL(r)
					instanceValue.Elem().FieldByIndex(baseUrlField.Index).SetString(baseURL)

				}

			}

		}

		// instanceValue := &result[0]

		return &result[0], nil
	})

	return *ret, err
}
