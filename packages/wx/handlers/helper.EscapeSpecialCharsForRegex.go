package handlers

import "strings"

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
