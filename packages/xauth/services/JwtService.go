package services

import (
	"encoding/json"
	"fmt"
	"time"
	"wx"
	"xauth/config"

	"github.com/golang-jwt/jwt/v5"
)

// Kết quả gồm cả access & refresh
type JWTPair struct {
	AccessToken  string
	RefreshToken string
}

type JWTService struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func (jwt *JWTService) New(configService *wx.Global[config.ConfigService]) error {
	var err error
	cfg, err := configService.Ins()
	if err != nil {
		return err
	}
	jwt.secret = []byte(cfg.Get().JwtToken.Secret)

	jwt.accessTTL, err = cfg.Get().JwtToken.TTL.GetAccess()
	if err != nil {
		return err
	}
	jwt.refreshTTL, err = cfg.Get().JwtToken.TTL.GetRefresh()
	if err != nil {
		return err
	}
	return nil
}
func (jwt *JWTService) GetSecret() []byte {
	return jwt.secret
}
func (jwt *JWTService) GetAccessTTL() time.Duration {
	return jwt.accessTTL
}
func (jwt *JWTService) GetRefreshTTL() time.Duration {
	return jwt.refreshTTL
}

func (jwtService *JWTService) Generate(data any) (*JWTPair, error) {
	// convert struct -> map
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	// ----- Access Token -----
	accessClaims := jwt.MapClaims{}
	for k, v := range m {
		accessClaims[k] = v
	}
	accessClaims["exp"] = time.Now().Add(jwtService.accessTTL).Unix()
	accessClaims["iat"] = time.Now().Unix()
	accessClaims["typ"] = "access"

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessStr, err := accessToken.SignedString(jwtService.secret)
	if err != nil {
		return nil, err
	}

	// ----- Refresh Token -----
	refreshClaims := jwt.MapClaims{
		"exp": time.Now().Add(jwtService.refreshTTL).Unix(),
		"iat": time.Now().Unix(),
		"typ": "refresh",
		"sub": m["user_id"], // thường chỉ giữ user_id cho gọn
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshStr, err := refreshToken.SignedString(jwtService.secret)
	if err != nil {
		return nil, err
	}

	return &JWTPair{
		AccessToken:  accessStr,
		RefreshToken: refreshStr,
	}, nil
}

// Hàm sinh access + refresh token từ struct

// Parse JWT về struct
func ParseToStruct[T any](secret []byte, tokenStr string) (*T, error) {
	parsedToken, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		// convert claims -> json -> struct
		b, err := json.Marshal(claims)
		if err != nil {
			return nil, err
		}
		var obj T
		if err := json.Unmarshal(b, &obj); err != nil {
			return nil, err
		}
		return &obj, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// ----------------- Demo -----------------
type Profile struct {
	UserID   int      `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	Email    string   `json:"email"`
}
