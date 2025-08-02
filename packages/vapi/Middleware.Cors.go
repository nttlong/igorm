package vapi

import "net/http"

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Thiết lập các header CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// Xử lý preflight request (OPTIONS)
		// Trình duyệt sẽ gửi request OPTIONS trước các request phức tạp
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Chuyển tiếp request đến handler tiếp theo
		next.ServeHTTP(w, r)
	})
}
func (server *HtttpServer) Cors() {
	oldHandler := server.handler
	if server.handler == nil {
		oldHandler = server.mux
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Thiết lập các header CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// Xử lý preflight request (OPTIONS)
		// Trình duyệt sẽ gửi request OPTIONS trước các request phức tạp
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		oldHandler.ServeHTTP(w, r)

	})
	server.handler = handler
}
