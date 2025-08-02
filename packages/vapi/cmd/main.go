package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"vapi"

	"github.com/golang-jwt/jwt/v5"
)

type TestController struct {
	HttpContext vapi.HttpContext
}

type Data struct {
	Code string
	Name string
}

func (c *TestController) Hello_Get() interface{} {
	return "Hello"
}
func (c *TestController) Test_Post(data Data, User vapi.UserClaims) Data {
	return data
}

type OAuth2PasswordImpl struct{}

var secretKey = []byte("my_super_secret_key")

func (o *OAuth2PasswordImpl) OnGetAccessToken(w http.ResponseWriter, r *http.Request) {
	// Chỉ chấp nhận phương thức POST
	if r.Method != http.MethodPost {
		http.Error(w, "Chỉ chấp nhận phương thức POST", http.StatusMethodNotAllowed)
		return
	}

	// Phân tích form data từ request
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Không thể phân tích form data", http.StatusBadRequest)
		return
	}

	// Lấy username và password từ form
	username := r.FormValue("username")
	password := r.FormValue("password")

	// 1. Giả lập quá trình xác thực
	// Trong thực tế, bạn sẽ truy vấn database hoặc một hệ thống xác thực khác.
	if username == "admin" && password == "password" {
		// 2. Nếu xác thực thành công, tạo một token (JWT)
		userID := "user-123"
		userRoles := []string{"admin", "member"}

		claims := vapi.UserClaims{
			UserID:   userID,
			Username: username,
			Roles:    userRoles,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				Issuer:    "my-auth-service",
				Subject:   userID,
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(secretKey)
		if err != nil {
			http.Error(w, "Lỗi khi tạo token", http.StatusInternalServerError)
			return
		}

		// 3. Giả lập một Refresh Token
		refreshTokenString := fmt.Sprintf("fake_refresh_token_%d", time.Now().Unix())

		// 4. Trả về cấu trúc AccessTokenResponse cho client
		response := vapi.AccessTokenResponse{
			AccessToken:  tokenString,
			TokenType:    "Bearer",
			ExpiresIn:    3600, // 1 giờ
			RefreshToken: refreshTokenString,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)

	} else {
		// 5. Nếu xác thực thất bại, trả về lỗi Unauthorized.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid_credentials"})
	}
}
func (o *OAuth2PasswordImpl) OnValidateToken(token string) vapi.TokenValidationResponse {
	claims := &vapi.UserClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return vapi.TokenValidationResponse{
			IsValid: false,
		}
	}
	return vapi.TokenValidationResponse{
		IsValid: true,
	}

}
func main() {
	vapi.AddHandler(func() (*TestController, error) {
		return &TestController{}, nil
	})
	server := vapi.NewHtttpServer(
		8080,
		"0.0.0.0",
	)
	server.OnValidateToken(func(token string) vapi.TokenValidationResponse {

		return vapi.TokenValidationResponse{
			IsValid: true,
		}

	})
	server.OnGetAccessToken(func(w http.ResponseWriter, r *http.Request) {
		// Chỉ chấp nhận phương thức POST
		if r.Method != http.MethodPost {
			http.Error(w, "Chỉ chấp nhận phương thức POST", http.StatusMethodNotAllowed)
			return
		}

		// Phân tích form data từ request
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Không thể phân tích form data", http.StatusBadRequest)
			return
		}

		// Lấy username và password từ form
		username := r.FormValue("username")
		password := r.FormValue("password")

		// 1. Giả lập quá trình xác thực
		// Trong thực tế, bạn sẽ truy vấn database hoặc một hệ thống xác thực khác.
		if username == "admin" && password == "password" {
			// 2. Nếu xác thực thành công, tạo một token (JWT)
			userID := "user-123"
			userRoles := []string{"admin", "member"}

			claims := vapi.UserClaims{
				UserID:   userID,
				Username: username,
				Roles:    userRoles,
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
					Issuer:    "my-auth-service",
					Subject:   userID,
				},
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString(secretKey)
			if err != nil {
				http.Error(w, "Lỗi khi tạo token", http.StatusInternalServerError)
				return
			}

			// 3. Giả lập một Refresh Token
			refreshTokenString := fmt.Sprintf("fake_refresh_token_%d", time.Now().Unix())

			// 4. Trả về cấu trúc AccessTokenResponse cho client
			response := vapi.AccessTokenResponse{
				AccessToken:  tokenString,
				TokenType:    "Bearer",
				ExpiresIn:    3600, // 1 giờ
				RefreshToken: refreshTokenString,
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)

		} else {
			// 5. Nếu xác thực thất bại, trả về lỗi Unauthorized.
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid_credentials"})
		}
	})
	server.Swagger()
	server.Cors()
	server.Auth()

	server.Start()
}
func main1() {
	vapi.AddHandler(func() (*TestController, error) {
		return &TestController{}, nil
	})
	// Tạo một file swagger.json giả để ví dụ hoạt động

	// Tạo một ServeMux mới. Bạn có thể sử dụng ServeMux có sẵn của ứng dụng.
	mux := http.NewServeMux()

	// Tích hợp middleware phục vụ Swagger UI
	vapi.ServeSwaggerUI(mux, &OAuth2PasswordImpl{})

	// In ra thông điệp và khởi động server

	vapi.Start(&OAuth2PasswordImpl{}, mux, "0.0.0.0", "8080")
}
