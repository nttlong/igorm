package vapi

import "net/http"

var Cors = func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// Cho phép tất cả origin (cẩn thận với sản phẩm thật!)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Nếu là preflight request (OPTIONS), chỉ phản hồi 200
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Gọi tiếp handler chính
	next.ServeHTTP(w, r)
}

// func (server *HtttpServer) Cors() {
// 	oldHandler := server.handler
// 	if server.handler == nil {
// 		oldHandler = server.mux
// 	}

// 	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Thiết lập các header CORS
// 		w.Header().Set("Access-Control-Allow-Origin", "*")
// 		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
// 		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

// 		// Xử lý preflight request (OPTIONS)
// 		// Trình duyệt sẽ gửi request OPTIONS trước các request phức tạp
// 		if r.Method == "OPTIONS" {
// 			w.WriteHeader(http.StatusOK)
// 			return
// 		}
// 		oldHandler.ServeHTTP(w, r)

// 	})
// 	server.Progress = append(server.Progress, handler)
// }
