package vapi

import (
	"fmt"
	"reflect"
	"strings"
)

func AddController[T any](initHanlder func() (*T, error)) error {

	swaggerData, err := loadSwaggerInfo()
	if err != nil {
		return err
	}

	typ := reflect.TypeFor[T]()
	fmt.Println(typ.String())
	typePrt := reflect.TypeFor[*T]()
	lst, err := utilsInstance.ParseMethods(typePrt)
	if err != nil {
		return err
	}
	utilsInstance.handler = append(utilsInstance.handler, lst...)
	utilsInstance.SortHandlers()

	for _, v := range lst {
		// ret := &PathItem{}
		pathItem := &PathItem{}

		retVal := reflect.ValueOf(pathItem).Elem()

		// 3. Lấy Operation và thông tin tham số từ hàm getOperation

		// 4. Lấy reflect.Type của PathItem (làm việc với kiểu)
		retTyp := reflect.TypeOf(&PathItem{}).Elem()
		op := utilsInstance.GetOperation(v)
		// 5. Tìm trường trong PathItem có tên tương ứng
		if field, ok := retTyp.FieldByNameFunc(func(s string) bool {
			return strings.EqualFold(s, v.HttpMethod)
		}); ok {
			// 6. Lấy reflect.Value của trường đó từ retVal (đã có thể ghi)
			fieldVal := retVal.FieldByIndex(field.Index)

			// 7. Kiểm tra tính hợp lệ và khả năng ghi
			if fieldVal.CanSet() {
				// 8. Gán con trỏ Operation vào trường
				fieldVal.Set(reflect.ValueOf(op))
			}
		}
		swaggerData.Paths[v.Url] = *pathItem

	}

	return nil

}
