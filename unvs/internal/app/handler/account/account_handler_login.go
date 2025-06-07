package account

import (
	"context"
	"errors"
	"fmt"
	"time"
	"unvs/internal/app/service/account"
	"unvs/internal/model/auth"

	_ "unvs/internal/config"
	config "unvs/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

// generateJWTToken generates a JWT token for a given user
func generateJWTToken(user *auth.User) (string, error) {
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
	signedString, err := token.SignedString(config.GetJWTSecret())
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return signedString, nil
}

// performLoginRequest là một helper method chứa logic chung cho việc xử lý đăng nhập.
// Nó nhận username (có thể là email) và password, sau đó gọi service và tạo token.
func (h *AccountHandler) performLoginRequest(ctx context.Context, username, password string) (*LoginResponse, error) {
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
	token, err := generateJWTToken(user)
	if err != nil {
		return nil, fmt.Errorf("TOKEN_GENERATION_FAILED: Không thể tạo token xác thực: %w", err)
	}

	// Trả về token và thông tin người dùng (không bao gồm mật khẩu hash)
	return &LoginResponse{
		Message: "Đăng nhập thành công",
		Token:   token,
		User:    user,
	}, nil
}
