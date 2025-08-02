package vapi

import "github.com/golang-jwt/jwt/v5"

// 2. Cấu trúc Payload (Claims) của Access Token (JWT).
// Đây là **cấu trúc quan trọng nhất** để phía server xác minh token.
// Nó chứa thông tin của người dùng và các quyền hạn (permissions).
// Cấu trúc này không bao giờ được trả về trực tiếp cho client.
type UserClaims struct {
	// Các trường chuẩn của JWT
	jwt.RegisteredClaims

	// Thông tin người dùng tùy chỉnh
	UserID   string   `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"` // Ví dụ: ["admin", "editor"]
}
