package wx

import "strings"

func (h *Handler) GetAbsRootUri() string {
	if h.rootAbsUrl != "" {
		return h.rootAbsUrl
	}
	scheme := h.GetScheme()

	h.rootAbsUrl = scheme + "://" + h.Req.Host
	return h.rootAbsUrl
}
func (h *Handler) GetScheme() string {
	if h.schema != "" {
		return h.schema
	}
	h.schema = "http"

	// Trường hợp Go server trực tiếp nhận HTTPS
	if h.Req.TLS != nil {
		return "https"
	}

	// Trường hợp có reverse proxy thêm header
	if proto := h.Req.Header.Get("X-Forwarded-Proto"); proto != "" {
		h.schema = strings.ToLower(proto)
	}

	if forwarded := h.Req.Header.Get("Forwarded"); forwarded != "" {
		// Ví dụ: "for=192.0.2.43; proto=https; by=203.0.113.43"
		for _, part := range strings.Split(forwarded, ";") {
			part = strings.TrimSpace(part)
			if strings.HasPrefix(strings.ToLower(part), "proto=") {
				h.schema = strings.TrimPrefix(part, "proto=")
			}
		}
	}

	// Mặc định là http
	return h.schema
}
