package services

type OAuth2Token struct {
	AccessToken  string  `json:"access_token"`
	TokenType    string  `json:"token_type"`
	ExpiresIn    int64   `json:"expires_in"` // Thời gian sống của token tính bằng giây
	Scope        string  `json:"scope"`
	RefreshToken string  `json:"refresh_token"`
	Message      string  `json:"message,omitempty"` // Thêm message nếu bạn muốn giữ lại
	RoleId       string  `json:"roleId,omitempty"`  // Thêm role nếu bạn muốn giữ lại
	UserId       string  `json:"userId,omitempty"`  // Thêm userID nếu bạn muốn giữ lại
	Username     string  `json:"username,omitempty"`
	Email        *string `json:"email,omitempty"`
}
