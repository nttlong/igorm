package vapi

// 1. Cấu trúc dữ liệu trả về cho client sau khi đăng nhập thành công.
// Cấu trúc này không đủ để xác minh, nó chỉ là thông tin cho client.
type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`              // Thời gian token hết hạn (đơn vị giây)
	RefreshToken string `json:"refresh_token,omitempty"` // Thường là tùy chọn
}
