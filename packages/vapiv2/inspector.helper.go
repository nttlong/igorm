package vapi

import (
	"reflect"
	"regexp"
	"strings"
)

func (h *helperType) ToKebabCase(s string) string {
	key := s + "/helperType/ToKebabCase"
	ret, _ := OnceCall(key, func() (*string, error) {
		// Khớp các chữ cái viết hoa và thêm dấu gạch ngang trước đó.
		// Ví dụ: MyMethod -> -My-Method
		re := regexp.MustCompile("([A-Z])")
		snake := re.ReplaceAllString(s, "-$1")

		// Chuyển toàn bộ chuỗi sang chữ thường và loại bỏ dấu gạch ngang ở đầu nếu có.
		// Ví dụ: -My-Method -> -my-method -> my-method
		ret := strings.ToLower(strings.TrimPrefix(snake, "-"))
		return &ret, nil
	})
	return *ret
}

func (h *helperType) EscapeSpecialCharsForRegex(s string) string {
	ret := ""
	for _, c := range s {
		if strings.Contains(h.SpecialCharForRegex, string(c)) {
			ret += "\\"
		}
		ret += string(c)
	}
	return ret
}
func GetMethodByName[T any](name string) *reflect.Method {

	t := reflect.TypeFor[*T]()
	key := t.Elem().String() + "/GetMethodByName/" + name
	ret, _ := OnceCall(key, func() (*reflect.Method, error) {

		for i := 0; i < t.NumMethod(); i++ {
			if t.Method(i).Name == name {
				ret := t.Method(i)
				return &ret, nil
			}
		}
		return nil, nil
	})
	return ret
}
