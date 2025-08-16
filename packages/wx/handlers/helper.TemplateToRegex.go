package handlers

import (
	"strings"
	"wx/internal"
)

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
