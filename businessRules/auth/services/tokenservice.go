package services

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	authErr "unvs.br.auth/errors"
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
type TokenService struct {
	CacheService
	JwtSecret []byte
}

func (p *TokenService) getPath() string {
	typ := reflect.TypeOf(*p)
	return typ.PkgPath() + "/" + typ.Name()
}
func (p *TokenService) validate() string {
	typ := reflect.TypeOf(*p)
	if p.Cache != nil {

		panic(fmt.Sprintf("\n%s.PasswordService.Cache is  nil\n ", typ.Name()))
	}
	if p.Context == nil {
		panic(fmt.Sprintf("\n%s.PasswordService.Context is  nil\n ", typ.Name()))
	}
	return typ.PkgPath() + "/" + typ.Name()
}

func (s *TokenService) ValidateAccessToken(accessToken string) (*OAuth2Token, error) {

	path := s.getPath()
	if accessToken == "" {
		e := &authErr.AuthError{
			Code:    authErr.ErrInvalidToken,
			Message: "access token is missing",
		}
		return nil, e
	}
	cacheKey := fmt.Sprintf("%s:%s", path, accessToken)
	ret := OAuth2Token{}
	if s.Cache.Get(s.Context, cacheKey, &ret) {
		return &ret, nil
	}
	accessTokenValidate := accessToken
	if strings.Contains(accessToken, "Bearer") {
		accessTokenValidate = strings.TrimPrefix(accessToken, "Bearer ")
	}
	tokenInfo, err := s.DecodeAccessToken(accessTokenValidate)
	if err != nil {
		if auErr, ok := err.(*authErr.AuthError); ok {
			return nil, auErr
		}
		return nil, err
	}

	s.Cache.Set(s.Context, cacheKey, tokenInfo, time.Duration(tokenInfo.ExpiresIn)*time.Second)

	return tokenInfo, nil
}

// decodeAccessToken giải mã accessToken và trả về OAuth2Token
func (s *TokenService) DecodeAccessToken(accessToken string) (*OAuth2Token, error) {
	// Parse token với claims

	token, err := jwt.ParseWithClaims(accessToken, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Kiểm tra thuật toán ký
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("unexpected signing method: %v", token.Header["alg"])
			return nil, &authErr.AuthError{
				Code:    authErr.ErrInvalidToken,
				Message: "invalid token",
			}
			//fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.JwtSecret, nil
	})
	if err != nil {
		if strings.Contains(err.Error(), ": token is expired") {
			return nil, &authErr.AuthError{
				Code:    authErr.ErrTokenExpired,
				Message: "token has expired",
			}
		}
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Kiểm tra token có hợp lệ không
	if !token.Valid {
		return nil, &authErr.AuthError{
			Code:    authErr.ErrInvalidToken,
			Message: "invalid token",
		}
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
		return nil, &authErr.AuthError{
			Code:    authErr.ErrInvalidToken,
			Message: "token has expired",
		}
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

// generateToken tạo một OAuth2Token với JWT và refresh token
func (s *TokenService) GenerateToken(userID string, role string) (*OAuth2Token, error) {
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
	accessToken, err := token.SignedString(s.JwtSecret)
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
