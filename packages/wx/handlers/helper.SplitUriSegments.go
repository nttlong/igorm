package handlers

// splitUriSegments splits the URI string by '/', ignoring empty segments.
func (h *helperType) SplitUriSegments(uri string) []string {
	var segments []string
	start := 0

	for i := 0; i < len(uri); i++ {
		if uri[i] == '/' {
			if start < i {
				segments = append(segments, uri[start:i])
			}
			start = i + 1
		}
	}
	// append the last segment if any
	if start < len(uri) {
		segments = append(segments, uri[start:])
	}
	return segments
}
