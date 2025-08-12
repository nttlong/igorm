package vapi

import (
	"fmt"
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
func GetUriOfHandler[T any](server *HtttpServer, methodName string) (string, error) {
	mt := GetMethodByName[T](methodName)
	if mt == nil {
		return "", fmt.Errorf("%s of %T was not found", methodName, *new(T))
	}
	mtInfo, err := inspector.helper.GetHandlerInfo(*mt)
	if err != nil {
		return "", fmt.Errorf("%s of %T cause  error %s", methodName, *new(T), err.Error())
	}
	if mtInfo == nil {
		return "", fmt.Errorf("%s of %T is not HttpMethod", methodName, *new(T))
	}
	if mtInfo.Uri != "" && mtInfo.Uri[0] == '/' {
		return mtInfo.Uri, nil
	}
	return server.BaseUrl + "/" + mtInfo.Uri, nil

}
