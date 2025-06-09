package oauth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"unvs/internal/app/service/account"
	"unvs/internal/config"
	"unvs/internal/model/auth"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

// LoginByFormSubmit xử lý yêu cầu HTTP POST để đăng nhập người dùng chỉ bằng form data.
// @Summary Đăng nhập người dùng bằng Form Submit (Username/Password)
// @Description Xác thực thông tin đăng nhập từ form data (username và mật khẩu) và trả về JWT token nếu thành công.
// @Tags Accounts
// @Accept x-www-form-urlencoded
// @Produce json
// @Param tenant formData string true "tenant code"
// @Param username formData string true "Tên người dùng (không phải email)"
// @Param password formData string true "Mật khẩu"
// @Param grant_type formData string false "Kiểu cấp quyền (thường là 'password' cho OAuth2, tùy chọn)"
// @Success 200 {object} OAuth2TokenResponse "Đăng nhập thành công, trả về JWT token"
// @Failure 400 {object} ErrorResponse "Yêu cầu không hợp lệ (validation errors)"
// @Failure 401 {object} ErrorResponse "Thông tin đăng nhập không hợp lệ"
// @Failure 500 {object} ErrorResponse "Lỗi nội bộ server"
// @Router /oauth/token [post]
func (h *OAuthHandler) Token(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	grantType := c.FormValue("grant_type") // Lấy grant_type nếu có

	// Có thể thêm validation cho grant_type nếu nó là bắt buộc cho luồng OAuth2 cụ thể
	if grantType != "" && grantType != "password" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "UNSUPPORTED_GRANT_TYPE",
			Message: "Chỉ hỗ trợ grant_type 'password'",
		})
	}

	// Gọi helper method để xử lý logic đăng nhập chính
	response, err := h.performLoginRequest(c.Request().Context(), username, password)
	if err != nil {
		// Phân tích lỗi từ performLoginRequest
		if strings.HasPrefix(err.Error(), "MISSING_CREDENTIALS") {
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Code:    "MISSING_CREDENTIALS",
				Message: strings.TrimPrefix(err.Error(), "MISSING_CREDENTIALS: "),
			})
		} else if strings.HasPrefix(err.Error(), "INVALID_CREDENTIALS") {
			return c.JSON(http.StatusUnauthorized, ErrorResponse{
				Code:    "INVALID_CREDENTIALS",
				Message: strings.TrimPrefix(err.Error(), "INVALID_CREDENTIALS: "),
			})
		} else if strings.HasPrefix(err.Error(), "TOKEN_GENERATION_FAILED") {
			return c.JSON(http.StatusInternalServerError, ErrorResponse{
				Code:    "TOKEN_GENERATION_FAILED",
				Message: strings.TrimPrefix(err.Error(), "TOKEN_GENERATION_FAILED: "),
			})
		} else { // AUTHENTICATION_ERROR hoặc lỗi không xác định khác
			return c.JSON(http.StatusInternalServerError, ErrorResponse{
				Code:    "AUTHENTICATION_ERROR",
				Message: "Lỗi xác thực không xác định",
			})
		}
	}

	return c.JSON(http.StatusOK, response)
}

// performLoginRequest là một helper method chứa logic chung cho việc xử lý đăng nhập.
// Nó nhận username (có thể là email) và password, sau đó gọi service và tạo token.
func (h *OAuthHandler) performLoginRequest(ctx context.Context, username, password string) (*OAuth2TokenResponse, error) {
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

// generateJWTToken (giữ nguyên hoặc sửa nhẹ để trả về expirationTime)
func generateJWTToken(user *auth.User) (string, time.Time, error) { // Thêm time.Time để trả về thời gian hết hạn
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
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}
	return signedString, expirationTime, nil
}
