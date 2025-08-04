package vapi

import (
	"encoding/json"
	"reflect"
	"strings"
)

type inputParamInfo struct {
	dataParamIndex int
	userParamIndex int
}

func getUserClaimsParamIndex(method reflect.Method) inputParamInfo {
	ret := inputParamInfo{
		dataParamIndex: -1,
		userParamIndex: -1,
	}
	for i := 1; i < method.Type.NumIn(); i++ {
		inputType := method.Type.In(i)
		if inputType.Kind() == reflect.Ptr {
			inputType = inputType.Elem()
		}
		if inputType == reflect.TypeOf(UserClaims{}) {
			ret.userParamIndex = i
		} else {
			ret.dataParamIndex = i

		}

	}
	return ret
}

func getParameters(method reflect.Method) ([]Parameter, inputParamInfo) {
	paramInfo := getUserClaimsParamIndex(method)

	if method.Type.NumIn() <= 1 { // only have reciver

		return []Parameter{}, paramInfo
	}

	ret := []Parameter{}
	typ := method.Type.In(1)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	sampleData := reflect.New(typ).Elem().Interface()
	jsonSampleData, err := json.Marshal(sampleData)
	if err != nil {
		return []Parameter{}, paramInfo
	}

	if typ.Kind() != reflect.Struct {
		return []Parameter{}, paramInfo
	}
	// for i := 0; i < typ.NumField(); i++ {
	// 	param := Parameter{
	// 		Type:        "string",
	// 		Description: "",
	// 		Name:        typ.Field(i).Name,
	// 		In:          "body",
	// 		Required:    false,
	// 	}
	// 	ret = append(ret, param)
	// }
	ret = append(ret, Parameter{
		Type:        "object",
		Description: string(jsonSampleData),
		Name:        typ.Name(),
		In:          "body",
	})
	return ret, paramInfo

}

func getOperation(method reflect.Method) (*Operation, inputParamInfo) {
	Parameters, paramInfo := getParameters(method)
	ret := Operation{
		Tags:       []string{},
		Summary:    "",
		Consumes:   []string{"application/json"},
		Parameters: Parameters,
		Responses: map[string]Response{
			"200": {
				Description: "",
				Schema:      nil,
			},
		},

		Produces: []string{"application/json"},
	}
	if paramInfo.userParamIndex != -1 {

		ret.Security = []map[string][]string{
			{
				"OAuth2Password": {},
			},
		}
	}

	return &ret, paramInfo
}

func getPathItem(methodName string, method reflect.Method, receiverType reflect.Type) (PathItem, inputParamInfo) {
	// 1. Tạo một con trỏ tới PathItem
	ret := &PathItem{}

	// 2. Lấy reflect.Value của con trỏ, sau đó Elem() để lấy giá trị có thể ghi
	retVal := reflect.ValueOf(ret).Elem()

	// 3. Lấy Operation và thông tin tham số từ hàm getOperation
	operation, paramInfo := getOperation(method)

	// 4. Lấy reflect.Type của PathItem (làm việc với kiểu)
	typ := reflect.TypeOf(&PathItem{}).Elem()

	// 5. Tìm trường trong PathItem có tên tương ứng
	if field, ok := typ.FieldByNameFunc(func(s string) bool {
		return strings.EqualFold(s, methodName)
	}); ok {
		// 6. Lấy reflect.Value của trường đó từ retVal (đã có thể ghi)
		fieldVal := retVal.FieldByIndex(field.Index)

		// 7. Kiểm tra tính hợp lệ và khả năng ghi
		if fieldVal.CanSet() {
			// 8. Gán con trỏ Operation vào trường
			fieldVal.Set(reflect.ValueOf(operation))
		}
	}

	// 9. Trả về giá trị struct, không phải con trỏ.
	return *ret, paramInfo
}

func getHttpContextFieldIndex(receiverType reflect.Type) []int {
	for i := 0; i < receiverType.NumField(); i++ {
		fieldType := receiverType.Field(i).Type
		if fieldType == reflect.TypeOf(HttpContext{}) {
			return receiverType.Field(i).Index

		} else if receiverType.Field(i).Anonymous {
			return getHttpContextFieldIndex(fieldType)
		}

	}
	return nil
}
