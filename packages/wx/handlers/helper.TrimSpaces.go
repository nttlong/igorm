package handlers

// trimSpaces trims leading and trailing spaces from a string.
func (h *helperType) TrimSpaces(s string) string {
	start, end := 0, len(s)-1
	for start <= end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end >= start && (s[end] == ' ' || s[end] == '\t') {
		end--
	}
	if start > end {
		return ""
	}
	return s[start : end+1]
}
