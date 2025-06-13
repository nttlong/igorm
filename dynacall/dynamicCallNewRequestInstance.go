package dynacall

import (
	"reflect"
	"strings"
)

type RequestType struct {
	Data   interface{}
	fields map[string]reflect.StructField
}

func (r *RequestType) GetString(key string) string {
	field, found := r.fields[strings.ToLower(key)]
	if found {
		return reflect.ValueOf(r.Data).Elem().FieldByName(field.Name).String()
	}
	return ""

}
func (r *RequestType) Get(key string) interface{} {
	field, found := r.fields[strings.ToLower(key)]
	if found {
		return reflect.ValueOf(r.Data).Elem().FieldByName(field.Name).Interface()
	}
	return nil

}
func NewRequestInstance(callerPath string, instanceType reflect.Type) (*RequestType, error) {
	ret := &RequestType{
		fields: make(map[string]reflect.StructField),
	}
	inputType, err := GetInputTypeOfCallerPath(callerPath)
	if err != nil {
		return nil, err
	}
	fields := make([]reflect.StructField, instanceType.NumField())
	for i := 0; i < instanceType.NumField(); i++ {

		if instanceType.Field(i).Name != "Args" {
			fields[i] = instanceType.Field(i)
			ret.fields[strings.ToLower(fields[i].Name)] = fields[i]
		} else {
			if len(inputType) == 0 {
				continue
			} else if len(inputType) == 1 {
				if inputType[0].Kind() == reflect.Ptr {
					panic("not implemented")
				}
				if inputType[0].Kind() == reflect.Slice {
					inputTypeSlice := reflect.SliceOf(inputType[0])
					fields[i] = reflect.StructField{
						Name: "Args",
						Type: inputTypeSlice,
						Tag:  reflect.StructTag(""),
					}
				} else {
					fields[i] = reflect.StructField{
						Name: "Args",
						Type: inputType[0],
						Tag:  reflect.StructTag(""),
					}
				}
			} else {
				interfaceType := reflect.TypeOf((*interface{})(nil)).Elem()
				dynamicSliceOfTypeInterface := reflect.SliceOf(interfaceType)
				sliceValue := reflect.MakeSlice(dynamicSliceOfTypeInterface, 0, 0)
				for _, t := range inputType {
					if t.Kind() == reflect.Ptr {
						t = t.Elem()
					}

					sliceValue = reflect.Append(sliceValue, reflect.New(t).Elem())
				}
				fields[i] = reflect.StructField{
					Name: "Args",
					Type: dynamicSliceOfTypeInterface,
					Tag:  reflect.StructTag(""),
				}

			}
			ret.fields[strings.ToLower(fields[i].Name)] = fields[i]

		}

	}
	structType := reflect.StructOf(fields)

	// Tạo instance của struct
	structInstance := reflect.New(structType).Elem()
	ret.Data = structInstance.Addr().Interface()

	return ret, nil
	return nil, nil
}
