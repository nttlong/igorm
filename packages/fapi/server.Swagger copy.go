package fapi

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

func (server *HtttpServer) Swagger() {
	mux := server.mux
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		// Đọc file swagger.json từ thư mục hiện tại

		// Thiết lập header để trình duyệt hiểu đây là file JSON
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("pong"))
	})
	// 1. Phục vụ file swagger.json từ đường dẫn /swagger.json
	// Thư viện httpSwagger sẽ tìm file này để hiển thị.
	mux.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		// Đọc file swagger.json từ thư mục hiện tại
		data := CreateMockSwaggerJSON()

		// Thiết lập header để trình duyệt hiểu đây là file JSON
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})

	// 2. Phục vụ giao diện Swagger UI trên đường dẫn /swagger/
	// Thư viện httpSwagger.WrapHandler tự động tạo giao diện HTML.
	// Đường dẫn thứ hai "./swagger.json" là vị trí của file JSON mà UI sẽ hiển thị.
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

}
