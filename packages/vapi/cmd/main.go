package main

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"time"
	"vapi"

	"github.com/golang-jwt/jwt/v5"
)

type TestController struct {
}

type Data struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// func (c *TestController) Hello() string {
// 	return "Hello"
// }

// func (c *TestController) Update(data Data, ctx vapi.HttpContext) Data {
// 	return data
// }

type FileDownload struct {
	AccessKey string `api:"url:{accessKey},description:file name,required:true"`
	FileName  string `api:"url:{fileName}.mp4,description:file name,required:true"`
}
type FileUpload struct {
	Tenant string `api:"url:{tenant},description:tenant name`
	File   *multipart.FileHeader
	Data   Data
}

func (c *TestController) Update(data Data, ctx vapi.HttpContext, user vapi.UserClaims) {

}

func (c *TestController) Upload(data FileUpload, ctx vapi.HttpContext, user vapi.UserClaims) {

}

func (c *TestController) File(data FileDownload, ctx vapi.HttpContext, user vapi.UserClaims) {

}

type TenanantIfno struct {
	Name string `json:"name"`
	Code string `json:"code"`
}
type UpdateTenant struct {
	Tenant string `api:"url:{tenant},description:tenant name:"tenant"`
	Info   TenanantIfno
}

func (c *TestController) UpdateTenant(data UpdateTenant, ctx vapi.HttpContext, user vapi.UserClaims) {

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
	// vapi.AddHandler(func() (*TestController, error) {
	// 	return &TestController{}, nil
	// })
	vapi.AddController(func() (*TestController, error) {
		return &TestController{}, nil
	})
	server := vapi.NewHtttpServer(
		8080,
		"0.0.0.0",
	)
	server.Middleware(vapi.Cors)
	server.Swagger()
	server.Middleware(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		next(w, r)
	})
	server.Middleware(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		next(w, r)

	})
	//server.Cors()
	//server.Auth()

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
