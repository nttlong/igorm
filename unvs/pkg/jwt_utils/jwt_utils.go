// pkg/jwt_utils/jwt_utils.go
package jwt_utils

import (
	"errors"
	"fmt"
	"time"

	"unvs/internal/config" // Import config package để lấy secret

	"github.com/golang-jwt/jwt/v4"
)

// UserClaims defines the claims for the JWT token
type UserClaims struct {
	UserID   string `json:"userId"` // Sử dụng string cho UUID
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateJWTToken generates a JWT token for a given user.
func GenerateJWTToken(userID, username, email string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token hết hạn sau 24 giờ
	claims := &UserClaims{
		UserID:   userID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSecret := config.GetJWTSecret()
	signedString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return signedString, nil
}

// DecodeJWT decodes a JWT token string into UserClaims.
func DecodeJWT(tokenString string) (*UserClaims, error) {
	jwtSecret := config.GetJWTSecret()

	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("phương thức ký không hợp lệ: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("lỗi khi parse hoặc xác minh token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token không hợp lệ hoặc đã hết hạn")
	}

	return claims, nil
}

// ContextKey định nghĩa kiểu cho key trong context
type ContextKey string

// UserClaimsContextKey là key để lưu UserClaims trong context
const UserClaimsContextKey ContextKey = "userClaims"
