package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type OAuth2Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"` // Thời gian sống của token tính bằng giây
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
	Message      string `json:"message,omitempty"` // Thêm message nếu bạn muốn giữ lại
	Role         string `json:"role,omitempty"`    // Thêm role nếu bạn muốn giữ lại
	UserId       string `json:"userId,omitempty"`  // Thêm userID nếu bạn muốn giữ lại
}

func verifyPassword(password, hashedPassword string) error {
	// So sánh mật khẩu với hash
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// hashPasswordWithSalt băm mật khẩu với muối sử dụng bcrypt
func hashPasswordWithSalt(password string) (string, error) {
	// Chuyển mật khẩu thành []byte
	passwordBytes := []byte(password)

	// Tạo hash với bcrypt, sử dụng cost factor mặc định (10)
	// bcrypt tự động tạo muối ngẫu nhiên
	hash, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	// Chuyển hash thành chuỗi để trả về
	return string(hash), nil
}

// generateToken tạo một OAuth2Token với JWT và refresh token
func generateToken(jwtSecret []byte, userID string, role string) (*OAuth2Token, error) {
	// Thời gian sống của token (ví dụ: 1 giờ)
	tokenDuration := 1 * time.Hour
	expirationTime := time.Now().Add(tokenDuration).Unix()

	// Tạo claims cho JWT
	claims := jwt.MapClaims{
		"userId": userID,            // userID là subject
		"exp":    expirationTime,    // Thời gian hết hạn
		"iat":    time.Now().Unix(), // Thời gian phát hành
		"scope":  "read write",      // Scope mặc định
		"role":   role,              // Role mặc định
	}

	// Tạo token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign JWT: %w", err)
	}

	// Tạo refresh token ngẫu nhiên
	refreshTokenBytes := make([]byte, 32)
	_, err = rand.Read(refreshTokenBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	refreshToken := base64.URLEncoding.EncodeToString(refreshTokenBytes)

	// Tạo OAuth2Token
	oauthToken := &OAuth2Token{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(tokenDuration.Seconds()),
		Scope:        "read write",
		RefreshToken: refreshToken,
		Message:      "Token generated successfully",
		Role:         role,
		UserId:       userID,
	}

	return oauthToken, nil
}

// decodeAccessToken giải mã accessToken và trả về OAuth2Token
func DecodeAccessToken(jwtSecret []byte, accessToken string) (*OAuth2Token, error) {
	// Parse token với claims
	token, err := jwt.ParseWithClaims(accessToken, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Kiểm tra thuật toán ký
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Kiểm tra token có hợp lệ không
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Lấy claims
	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to parse claims")
	}

	// Trích xuất các trường từ claims
	userID, ok := (*claims)["userId"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid userId claim")
	}
	scope, _ := (*claims)["scope"].(string)
	role, _ := (*claims)["role"].(string)
	exp, ok := (*claims)["exp"].(float64)
	if !ok {
		return nil, fmt.Errorf("missing or invalid exp claim")
	}

	// Tính ExpiresIn (thời gian còn lại tính bằng giây)
	currentTime := time.Now().Unix()
	expiresIn := int64(exp) - currentTime
	if expiresIn < 0 {
		return nil, fmt.Errorf("token has expired")
	}

	// Tạo OAuth2Token
	oauthToken := &OAuth2Token{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
		Scope:        scope,
		RefreshToken: "", // RefreshToken không có trong accessToken
		Message:      "Token decoded successfully",
		Role:         role,
		UserId:       userID,
	}

	return oauthToken, nil
}
