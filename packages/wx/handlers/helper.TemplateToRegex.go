package handlers

import (
	"regexp"
	"strings"
	"wx/internal"
)

func (h *helperType) convertUrlToRegex(urlPattern string) string {
	// Bước 1: Thay thế các wildcard catch-all (*...) bằng .*
	//regexPattern := strings.ReplaceAll(urlPattern, "*", ".*")

	// Bước 2: Xử lý các tham số đường dẫn thông thường {name}
	// Ví dụ: {id} sẽ được chuyển thành ([^/]+) để khớp với bất kỳ ký tự nào ngoại trừ "/"
	// re := regexp.MustCompile(`{[^}]+}`)
	re := regexp.MustCompile(`\{[*][^{}]+\}`)
	reParam := regexp.MustCompile(`\{[^{}*]+\}`)
	// Thay thế tất cả các khớp với biểu thức (.*)
	// Lưu ý: Chúng ta dùng (.*) để bắt toàn bộ nội dung, bao gồm cả dấu gạch chéo
	regexPattern := re.ReplaceAllString(urlPattern, "(.*)")
	regexPattern = reParam.ReplaceAllString(regexPattern, "([^/]+)")

	return regexPattern
}

// templateToRegex chuyển URI template thành regex pattern string
// lấy các giá trị trong {}
func (h *helperType) TemplateToRegex(template string) string {
	key := template + "/helperType/TemplateToRegex"
	ret, _ := internal.OnceCall(key, func() (*string, error) {
		segments := strings.Split(template, "/")
		regexParts := []string{}
		paramCount := 0
		var escapeRegex = h.EscapeSpecialCharsForRegex
		for _, seg := range segments {
			if seg == "" {
				continue
			}

			var sb strings.Builder
			i := 0
			for i < len(seg) {
				start := strings.Index(seg[i:], "{")
				if start == -1 {
					// No more '{', escape remainder
					sb.WriteString(escapeRegex(seg[i:]))
					break
				}

				start += i
				end := strings.Index(seg[start:], "}")
				if end == -1 {
					// No closing brace, treat literally
					sb.WriteString(escapeRegex(seg[i:]))
					break
				}
				end += start

				// Escape static part before {
				if start > i {
					sb.WriteString(escapeRegex(seg[i:start]))
				}

				// Add capture group for parameter inside {}
				sb.WriteString(`([^/]+)`)
				paramCount++

				// Move index past "}"
				i = end + 1
			}

			regexParts = append(regexParts, sb.String())
		}

		// Join parts with '/'
		regexPattern := "^" + strings.Join(regexParts, "/") + "$"
		return &regexPattern, nil
	})
	return *ret
}
