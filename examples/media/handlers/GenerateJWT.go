package handlers

import (
	"encoding/json"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Hàm generic tạo JWT từ bất kỳ struct nào
func GenerateJWT[T any](data T, secret string, expire time.Duration) (string, error) {
	// 1. Convert struct sang map để nhét vào claims
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	claimsMap := make(map[string]any)
	if err := json.Unmarshal(jsonBytes, &claimsMap); err != nil {
		return "", err
	}

	// 2. Thêm exp (expire)
	claimsMap["exp"] = time.Now().Add(expire).Unix()

	// 3. Tạo token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claimsMap))

	// 4. Ký với secret
	return token.SignedString([]byte(secret))
}

type LoginService struct {
	SharedSecret  string
	DefaultExpire time.Duration
}

func (s *LoginService) New() error {
	s.SharedSecret = "mysupersecretkeythatismorethan32chars!"
	s.DefaultExpire = time.Hour * 24 // Mặc định là 24 giờ
	return nil
}
func (s *LoginService) GenerateJWT(data any, expire *time.Duration) (string, error) {
	// 1. Convert struct sang map để nhét vào claims
	expireValue := s.DefaultExpire
	if expire != nil {
		expireValue = *expire
	}
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	claimsMap := make(map[string]any)
	if err := json.Unmarshal(jsonBytes, &claimsMap); err != nil {
		return "", err
	}

	// 2. Thêm exp (expire)
	claimsMap["exp"] = time.Now().Add(expireValue).Unix()

	// 3. Tạo token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claimsMap))

	// 4. Ký với secret
	return token.SignedString([]byte(s.SharedSecret))
}
