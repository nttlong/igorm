package vapi

import (
	"reflect"
	"regexp"
	"strings"
)

func ToKebabCase(s string) string {
	// Khớp các chữ cái viết hoa và thêm dấu gạch ngang trước đó.
	// Ví dụ: MyMethod -> -My-Method
	re := regexp.MustCompile("([A-Z])")
	snake := re.ReplaceAllString(s, "-$1")

	// Chuyển toàn bộ chuỗi sang chữ thường và loại bỏ dấu gạch ngang ở đầu nếu có.
	// Ví dụ: -My-Method -> -my-method -> my-method
	return strings.ToLower(strings.TrimPrefix(snake, "-"))
}

const specialCharForRegex = "/\\?.$%^*-+"

func EscapeSpecialCharsForRegex(s string) string {
	ret := ""
	for _, c := range s {
		if strings.Contains(specialCharForRegex, string(c)) {
			ret += "\\"
		}
		ret += string(c)
	}
	return ret
}
func GetMethodByName[T any](name string) *reflect.Method {
	t := reflect.TypeFor[*T]()
	for i := 0; i < t.NumMethod(); i++ {
		if t.Method(i).Name == name {
			ret := t.Method(i)
			return &ret
		}
	}
	return nil
}
