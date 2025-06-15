package dynacall

import (
	"errors"
	"reflect"
	"strings"
)

// func (r *RequestType) GetString(key string) string {
// 	field, found := r.fields[strings.ToLower(key)]
// 	if found {
// 		return reflect.ValueOf(r.Data).Elem().FieldByName(field.Name).String()
// 	}
// 	return ""

// }
// func (r *RequestType) Get(key string) interface{} {
// 	field, found := r.fields[strings.ToLower(key)]
// 	if found {
// 		val := reflect.ValueOf(r.Data).Elem().FieldByName(field.Name)
// 		ptr := reflect.New(reflect.ValueOf(r.Data).Elem().FieldByName(field.Name).Type()) // Tạo một con trỏ mới đến kiểu của val (ví dụ: *main.MyData)
// 		ptr.Elem().Set(val)                                                               // Gán giá trị của val vào vùng nhớ mà ptr trỏ tới
// 		return ptr.Interface()
// 	}

// 	return nil

// }

func GetInputExampleCallerPath(callerPath string) ([]interface{}, error) {
	if !strings.Contains(callerPath, "@") {
		return nil, errors.New("callerPath is invalid")
	}
	callerEntry, found := callerCache.Load(strings.ToLower(callerPath))
	if !found {
		return nil, errors.New("callerEntry not found")
	}
	method := callerEntry.(CallerEntry).Method
	ret := []interface{}{}
	for i := 1; i < method.Type.NumIn(); i++ {
		inputType := method.Type.In(i)
		if inputType.Kind() == reflect.Ptr {
			inputType = inputType.Elem()
		}
		inputValue := reflect.New(inputType).Interface()
		ret = append(ret, inputValue)
	}
	return ret, nil
}
