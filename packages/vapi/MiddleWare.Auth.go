package vapi

import (
	"log"
	"net/http"
	"regexp"
	"strings"
)

func AuthMiddleware(next http.Handler, onValidateToken func(token string) TokenValidationResponse) http.Handler {
	if onValidateToken == nil {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Lấy giá trị của header "Authorization"
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Thiếu Authorization Header", http.StatusUnauthorized)
			return
		}

		// 2. Tách chuỗi token từ định dạng "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Authorization Header không đúng định dạng", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		// 3. Gọi hàm validateToken để xác minh token
		validationResult := onValidateToken(token)
		if !validationResult.IsValid {
			http.Error(w, "Token không hợp lệ hoặc đã hết hạn", http.StatusUnauthorized)
			return
		}

		// 4. Nếu token hợp lệ, chuyển request đến handler tiếp theo
		log.Printf("Token hợp lệ. UserID: %s, Username: %s", validationResult.UserID, validationResult.Username)
		next.ServeHTTP(w, r)
	})
}
func (server *HtttpServer) Auth() {
	if server.onGetAccessToken == nil {
		panic("onGetAccessToken is nil,please call server.OnGetAccessToken(func(w http.ResponseWriter, r *http.Request))")
	}
	if server.onValidateToken == nil {
		panic("onValidateToken is nil,please call server.OnValidateToken(func(token string) TokenValidationResponse)")

	}

	oldHandler := server.handler
	if server.handler == nil {
		oldHandler = server.mux
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if server.urlAuth == r.URL.Path {

			server.onGetAccessToken(w, r)
			return
		}
		if ret, ok := server.urlPassAuth[r.URL.Path]; ok {
			if ret {
				oldHandler.ServeHTTP(w, r)
				return
			}
		}
		for _, re := range server.publicUrl {
			if re.MatchString(r.URL.Path) {
				oldHandler.ServeHTTP(w, r)
				server.urlPassAuth[r.URL.Path] = true
				return
			}
		}

		// 1. Lấy giá trị của header "Authorization"
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Thiếu Authorization Header", http.StatusUnauthorized)
			return
		}

		// 2. Tách chuỗi token từ định dạng "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Authorization Header không đúng định dạng", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		// 3. Gọi hàm validateToken để xác minh token
		validationResult := server.onValidateToken(token)
		if !validationResult.IsValid {
			http.Error(w, "Token không hợp lệ hoặc đã hết hạn", http.StatusUnauthorized)
			return
		}

		// 4. Nếu token hợp lệ, chuyển request đến handler tiếp theo
		log.Printf("Token hợp lệ. UserID: %s, Username: %s", validationResult.UserID, validationResult.Username)
		oldHandler.ServeHTTP(w, r)
	})
	server.publicUrl = append(server.publicUrl, regexp.MustCompile("^/oauth/token$"))
	server.handler = handler
	server.urlAuth = "/oauth/token"
}
func (server *HtttpServer) OnValidateToken(onValidateToken func(token string) TokenValidationResponse) {
	server.onValidateToken = onValidateToken

}
