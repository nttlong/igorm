package oauth

// OAuthHandler là struct chứa dependency đến AccountService.
type OAuthHandler struct {
}

// ErrorResponse struct for consistent error messages
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// type UserClaims struct {
// 	UserID   string `json:"userId"` // Sử dụng string cho UUID
// 	Username string `json:"username"`
// 	Email    string `json:"email"`
// 	jwt.RegisteredClaims
// }

// // OAuth2TokenResponse là cấu trúc phản hồi chuẩn cho OAuth2 Password Flow
// type OAuth2TokenResponse struct {
// 	AccessToken string `json:"access_token"`
// 	TokenType   string `json:"token_type"`
// 	ExpiresIn   int64  `json:"expires_in"` // Thời gian sống của token tính bằng giây
// 	// Các trường khác có thể thêm vào nếu cần như refresh_token, scope, v.v.
// 	Message string     `json:"message,omitempty"` // Thêm message nếu bạn muốn giữ lại
// 	User    *auth.User `json:"user,omitempty"`    // Thêm thông tin user nếu bạn muốn giữ lại
// }
