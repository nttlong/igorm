package handlers

// extractNameInBraces extracts the trimmed content inside the first pair of braces '{}' in the segment.
// Returns empty string if no braces found.
func (h *helperType) ExtractNameInBraces(segment string) string {
	start := -1
	end := -1
	for i, ch := range segment {
		if ch == '{' && start == -1 {
			start = i
		} else if ch == '}' && start != -1 {
			end = i
			break
		}
	}
	if start != -1 && end != -1 && end > start+1 {
		name := segment[start+1 : end]
		return h.TrimSpaces(name)
	}
	return ""
}
