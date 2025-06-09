package accounts

import (
	"context"
	"errors"
	"fmt"
	"time"
	"unvs/internal/app/service/account"
	"unvs/internal/model/auth"

	_ "unvs/internal/config"
	config "unvs/internal/config"

	"github.com/golang-jwt/jwt/v4"
)

// OAuth2TokenResponse là cấu trúc phản hồi chuẩn cho OAuth2 Password Flow
type OAuth2TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"` // Thời gian sống của token tính bằng giây
	// Các trường khác có thể thêm vào nếu cần như refresh_token, scope, v.v.
	Message string     `json:"message,omitempty"` // Thêm message nếu bạn muốn giữ lại
	User    *auth.User `json:"user,omitempty"`    // Thêm thông tin user nếu bạn muốn giữ lại
}

// generateJWTToken tạo token JWT và trả về token string cùng thời gian hết hạn
func generateJWTToken(user *auth.User) (string, time.Time, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token hết hạn sau 24 giờ
	claims := &UserClaims{
		UserID:   user.UserId,
		Username: user.Username,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.UserId,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte(config.GetJWTSecret()))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}
	return signedString, expirationTime, nil
}

// CreateOAuth2TokenResponse tạo phản hồi chuẩn OAuth2 từ thông tin token và user
func CreateOAuth2TokenResponse(user *auth.User) (*OAuth2TokenResponse, error) {
	tokenString, expirationTime, err := generateJWTToken(user)
	if err != nil {
		return nil, err
	}

	// Tính thời gian sống còn lại của token bằng giây
	expiresIn := expirationTime.Unix() - time.Now().Unix()
	if expiresIn < 0 { // Đảm bảo không trả về giá trị âm nếu token đã hết hạn
		expiresIn = 0
	}

	return &OAuth2TokenResponse{
		AccessToken: tokenString,
		TokenType:   "Bearer",
		ExpiresIn:   expiresIn,
		Message:     "Đăng nhập thành công", // Giữ lại nếu muốn
		User:        user,                   // Giữ lại nếu muốn
	}, nil
}

// performLoginRequest là một helper method chứa logic chung cho việc xử lý đăng nhập.
// Nó nhận username (có thể là email) và password, sau đó gọi service và tạo token.
func (h *AccountHandler) performLoginRequest(ctx context.Context, username, password string) (*OAuth2TokenResponse, error) {
	// Validate input (có thể dùng validator riêng nếu phức tạp hơn)
	if username == "" || password == "" {
		return nil, errors.New("MISSING_CREDENTIALS: Username và mật khẩu không được để trống")
	}

	// Gọi phương thức xác thực từ Service Layer
	// Service sẽ tự quyết định xác thực bằng email hay username
	user, err := h.accountService.AuthenticateUser(ctx, username, password)
	if err != nil {
		// Dịch lỗi từ Service sang lỗi HTTP phù hợp
		switch {
		case errors.Is(err, account.ErrInvalidCredentials):
			return nil, errors.New("INVALID_CREDENTIALS: Username hoặc mật khẩu không đúng")
		case errors.Is(err, account.ErrUserNotFound): // Có thể không cần thiết vì ErrInvalidCredentials bao trùm
			return nil, errors.New("INVALID_CREDENTIALS: Username hoặc mật khẩu không đúng")
		default:
			// Các lỗi khác từ service (ví dụ: lỗi DB không mong muốn)
			return nil, fmt.Errorf("AUTHENTICATION_ERROR: Lỗi xác thực không xác định: %w", err)
		}
	}

	// Tạo JWT token sau khi xác thực thành công
	auth2Token, err := CreateOAuth2TokenResponse(user)

	if err != nil {
		return nil, fmt.Errorf("TOKEN_GENERATION_FAILED: Không thể tạo token xác thực: %w", err)
	}

	// Trả về token và thông tin người dùng (không bao gồm mật khẩu hash)
	return auth2Token, nil
}
