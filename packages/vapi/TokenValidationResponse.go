package vapi

// 3. Cấu trúc dữ liệu phản hồi cho các hệ thống khác.
// Khi một hệ thống khác cần xác minh token, nó sẽ gọi một endpoint
// của bạn và nhận lại cấu trúc này.
type TokenValidationResponse struct {
	IsValid  bool     `json:"is_valid"`
	UserID   string   `json:"user_id,omitempty"`
	Username string   `json:"username,omitempty"`
	Roles    []string `json:"roles,omitempty"`
	// Có thể thêm các trường khác như IsExpired, Error, v.v.
}
