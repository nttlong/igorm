package vapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	httpSwagger "github.com/swaggo/http-swagger"
)

// ServeSwaggerUI đóng gói logic để phục vụ file swagger.json và giao diện Swagger UI.
// Hàm này có thể được xem như một middleware và dễ dàng tích hợp vào bất kỳ ServeMux nào.
func ServeSwaggerUI(mux *http.ServeMux, auth OAuth2Password) {
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

	for k := range mapHandler {
		fmt.Println(k)
		mux.HandleFunc(k, func(w http.ResponseWriter, r *http.Request) {

			method, ok := mapHandler[k]
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			httpMethod := httpMethodMap[k]
			if r.Method != httpMethod {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			if httpMethod == "POST" {
				serverResolverPost(w, r, method, k)

			} else {

				instanceVal := mapInstanceInit[k]
				output := method.Func.Call([]reflect.Value{instanceVal})
				w.Header().Set("Content-Type", "application/json")
				outPutValues := make([]interface{}, len(output))
				for i, v := range output {
					outPutValues[i] = v.Interface()
				}

				jsonData, _ := json.Marshal(outPutValues[0])

				w.Write(jsonData)
			}

		})
	}

	// 2. Phục vụ giao diện Swagger UI trên đường dẫn /swagger/
	// Thư viện httpSwagger.WrapHandler tự động tạo giao diện HTML.
	// Đường dẫn thứ hai "./swagger.json" là vị trí của file JSON mà UI sẽ hiển thị.
	mux.Handle("/swagger/", httpSwagger.WrapHandler)
}

type OAuth2Password interface {
	OnGetAccessToken(w http.ResponseWriter, r *http.Request)
	OnValidateToken(token string) TokenValidationResponse
}
