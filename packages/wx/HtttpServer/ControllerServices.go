package htttpserver

import (
	"net"
	"net/http"
	"sync"
)

type ContetxService struct {
	BaseUrl string
}

var onceGetBaseUrl sync.Once

// GetBaseURL lấy scheme, hostname, port và base URL từ request HTTP
func getBaseURL(r *http.Request) (scheme, hostname, port, baseURL string) {
	getBaseURL := func(r *http.Request) (scheme, hostname, port, baseURL string) {
		// Xác định scheme

		scheme = "http"
		if r.TLS != nil {
			scheme = "https"
		}

		// Lấy host (có thể bao gồm port)
		host := r.Host

		// Tách hostname và port
		hostname, port, err := net.SplitHostPort(host)
		if err != nil {
			// Nếu không có port trong host
			hostname = host
			port = ""
		}

		// Ghép baseURL
		baseURL = scheme + "://" + hostname
		if port != "" && !((scheme == "http" && port == "80") || (scheme == "https" && port == "443")) {
			baseURL += ":" + port
		}
		return scheme, hostname, port, baseURL
	}
	onceGetBaseUrl.Do(func() {
		scheme, hostname, port, baseURL = getBaseURL(r)
	})
	return scheme, hostname, port, baseURL
}
