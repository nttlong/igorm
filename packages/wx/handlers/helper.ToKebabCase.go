package handlers

import (
	"regexp"
	"strings"
	"wx/internal"
)

func (h *helperType) ToKebabCase(s string) string {
	key := s + "/helperType/ToKebabCase"
	ret, _ := internal.OnceCall(key, func() (*string, error) {
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
