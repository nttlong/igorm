package htttpserver

import (
	"net/http"
	"reflect"
	"wx/internal"
)

var IsTypeDepen func(t reflect.Type, visited map[reflect.Type]bool) bool
var CreateDepen func(depenType reflect.Type) (*reflect.Value, error)
var FindNewMethod func(t reflect.Type) (*reflect.Method, error)
var ResolveNewMethod func(mt reflect.Method) (*reflect.Value, error)
var ResolveNewMethodWithReceiver func(retVale reflect.Value, nm reflect.Method) (*reflect.Value, error)

func (web *webHandlerRunnerType) ResolveReceiverValue(handler WebHandler, r *http.Request) (reflect.Value, error) {
	key := handler.ApiInfo.ReceiverType.String() + "/webHandlerRunnerType/ResolveReceiverValue"
	ret, err := internal.OnceCall(key, func() (*reflect.Value, error) {

		result := handler.InitFunc.Call([]reflect.Value{})
		if result[1].IsValid() && !result[1].IsNil() {
			return nil, result[1].Interface().(error)
		}
		instanceType := handler.ApiInfo.ReceiverType
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
				if parentBaseurlFieldType == reflect.TypeOf(ContetxService{}) {
					instanceValue := &result[0]
					_, _, _, baseURL := getBaseURL(r)
					instanceValue.Elem().FieldByIndex(baseUrlField.Index).SetString(baseURL)

				}

			}

		}

		for i := 0; i < handler.ApiInfo.ReceiverTypeElem.NumField(); i++ {
			field := handler.ApiInfo.ReceiverTypeElem.Field(i)
			fieldType := field.Type
			newOfReiverType, _ := FindNewMethod(fieldType)
			if newOfReiverType != nil {
				method := *newOfReiverType
				r, err := ResolveNewMethod(method)
				if err != nil {
					return nil, err
				}
				reciverType := result[0]
				if reciverType.Kind() == reflect.Ptr {
					reciverType = reciverType.Elem()
				}
				fieldSet := reciverType.Field(i)
				valueSet := *r

				if fieldSet.IsValid() {
					fieldSet.Set(valueSet.Elem())
				}

			}
			if IsTypeDepen(fieldType, map[reflect.Type]bool{}) {
				if fieldType.Kind() == reflect.Ptr {
					fieldType = fieldType.Elem()
				}
				dependVale := reflect.New(fieldType)
				reciverType := result[0]
				if reciverType.Kind() == reflect.Ptr {
					reciverType = reciverType.Elem()
				}

				appField := dependVale.Elem().FieldByName("Owner")
				if appField.IsValid() {
					if appField.Kind() == reflect.Ptr {
						appField.Set(result[0])
					} else {
						appField.Set(reflect.ValueOf(result[0].Interface()))
					}

				}
				fieldSet := reciverType.Field(i)
				if fieldSet.IsValid() {
					if field.Type.Kind() == reflect.Ptr {
						fieldSet.Set(dependVale)
					} else {
						fieldSet.Set(dependVale.Elem())
					}

				}
			}

		}
		newMethodOfReiverType, err := FindNewMethod(handler.ApiInfo.ReceiverType)
		if err != nil {
			return nil, err
		}
		if newMethodOfReiverType != nil {
			method := *newMethodOfReiverType
			args := make([]reflect.Value, method.Type.NumIn())
			args[0] = result[0]
			for j := 1; j < method.Type.NumIn(); j++ {
				argTyp := method.Type.In(j)
				if IsTypeDepen(argTyp, map[reflect.Type]bool{}) {
					if argTyp.Kind() == reflect.Ptr {
						argTyp = argTyp.Elem()
					}
					dependVale := reflect.New(argTyp)
					reciverType := result[0]
					if reciverType.Kind() == reflect.Ptr {
						reciverType = reciverType.Elem()
					}

					appField := dependVale.Elem().FieldByName("Owner")
					if appField.IsValid() {
						if appField.Kind() == reflect.Ptr {
							appField.Set(result[0])
						} else {
							appField.Set(reflect.ValueOf(result[0].Interface()))
						}

					}
					if method.Type.In(j).Kind() == reflect.Ptr {
						args[j] = dependVale
					} else {
						args[j] = dependVale.Elem()
					}

				} else {
					if argTyp.Kind() == reflect.Ptr {
						argTyp = argTyp.Elem()
					}
					mt, err := FindNewMethod(argTyp)
					if err != nil {
						return nil, err
					}
					argVal, err := ResolveNewMethod(*mt)
					if err != nil {
						return nil, err
					}
					if method.Type.In(j).Kind() == reflect.Ptr {
						args[j] = *argVal
					} else {
						args[j] = (*argVal).Elem()
					}
				}

			}

			ret := method.Func.Call(args)
			if len(ret) > 0 && !ret[0].IsNil() {
				return nil, ret[0].Interface().(error)
			}
		}

		// instanceValue := &result[0]

		return &result[0], nil
	})

	return *ret, err
}
