package caller

import (
	"errors"
	"time" // Để làm việc với thời gian hết hạn
	"unvs/views"

	"github.com/golang-jwt/jwt/v4" // Import thư viện JWT
)

// jwtDecodeInfo là struct custom claims của bạn.
// Nó chứa các claim riêng của ứng dụng bạn, và nhúng jwt.RegisteredClaims
// để xử lý các claim chuẩn của JWT.

// jwtDecode giải mã token JWT và trả về thông tin giải mã, với kiểm tra expire
func jwtDecode(tokenString string, jwtSecret string) (*views.JwtDecodeInfo, error) {
	// Kiểm tra token rỗng
	if tokenString == "" {
		return nil, errors.New("token string is empty")
	}

	// Tạo key để xác minh token
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		// Kiểm tra thuật toán ký (HMAC trong trường hợp này)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	}

	// Giải mã token
	token, err := jwt.ParseWithClaims(tokenString, &views.JwtDecodeInfo{}, keyFunc)
	if err != nil {
		// Kiểm tra lỗi cụ thể nếu token hết hạn
		var validationError *jwt.ValidationError
		if errors.As(err, &validationError) {
			if validationError.Inner != nil && errors.Is(validationError.Inner, jwt.ErrTokenExpired) {
				return nil, errors.New("token has expired")
			}
		}
		return nil, err
	}

	// Kiểm tra token hợp lệ
	if claims, ok := token.Claims.(*views.JwtDecodeInfo); ok && token.Valid {
		// Kiểm tra expire thủ công để đảm bảo
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			return nil, errors.New("token has expired")
		}
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
