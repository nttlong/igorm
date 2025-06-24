package services

import (
	"context"
	"dbx"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	cacher "unvs.core/cacher"
	_ "unvs.core/errors"
	coreErrors "unvs.core/errors"

	_ "unvs.core/errors"

	"github.com/golang-jwt/jwt/v4"
)

type TokenService struct {
	Cache   cacher.Cache
	Context context.Context

	TenantDb      *dbx.DBXTenant
	EncryptionKey string
	Language      string
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
		e := &coreErrors.CoreError{
			Code: coreErrors.Error_AccessDenied,
			Err:  errors.New("Access Deny"),
		}
		return nil, e
	}
	cacheKey := fmt.Sprintf("%s:%s;%s", s.TenantDb.TenantDbName, path, accessToken)
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

		return nil, err
	}
	if tokenInfo == nil {
		return nil, &coreErrors.CoreError{
			Code:    coreErrors.Error_AccessDenied,
			Message: "Access Deny",
			Err:     errors.New("Access Deny"),
		}
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
			return nil, &coreErrors.CoreError{
				Code:    coreErrors.Error_AccessDenied,
				Message: "invalid token",
			}
			//fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		jwtSecret, err := s.GetJwtSecret()
		if err != nil {
			return nil, err
		}
		return jwtSecret, nil
	})
	if err != nil {
		if strings.Contains(err.Error(), ": token is expired") {
			return nil, &coreErrors.CoreError{
				Code:    coreErrors.Error_TokenExpired,
				Message: "invalid token",
				Err:     errors.New("token has expired"),
			}
		}
		if strings.Contains(err.Error(), ": signature is invalid") {
			return nil, &coreErrors.CoreError{
				Code:    coreErrors.Error_AccessDenied,
				Message: "Access Deny",
				Err:     errors.New("signature is invalid"),
			}
		}
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Kiểm tra token có hợp lệ không
	if !token.Valid {
		return nil, &coreErrors.CoreError{
			Code:    coreErrors.Error_AccessDenied,
			Message: "Access Deny",
			Err:     errors.New("token is invalid"),
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
	roleId, _ := (*claims)["roleId"].(string)
	exp, ok := (*claims)["exp"].(float64)
	if !ok {
		return nil, fmt.Errorf("missing or invalid exp claim")
	}

	// Tính ExpiresIn (thời gian còn lại tính bằng giây)
	currentTime := time.Now().Unix()
	expiresIn := int64(exp) - currentTime
	if expiresIn < 0 {
		return nil, &coreErrors.CoreError{
			Code:    coreErrors.Error_TokenExpired,
			Message: "Token has expired",
			Err:     errors.New("token has expired"),
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
		RoleId:       roleId,
		UserId:       userID,
	}

	return oauthToken, nil
}

type GenerateTokenParams struct {
	UserId   string
	RoleId   string
	Username string
	Email    *string
}

// generateToken tạo một OAuth2Token với JWT và refresh token
func (s *TokenService) GenerateToken(data GenerateTokenParams) (*OAuth2Token, error) {
	// Thời gian sống của token (ví dụ: 1 giờ)
	tokenDuration := 1 * time.Hour
	expirationTime := time.Now().Add(tokenDuration).Unix()

	// Tạo claims cho JWT
	claims := jwt.MapClaims{
		"userId":   data.UserId,       // userID là subject
		"exp":      expirationTime,    // Thời gian hết hạn
		"iat":      time.Now().Unix(), // Thời gian phát hành
		"scope":    "read write",      // Scope mặc định
		"role":     data.UserId,       // Role mặc định
		"username": data.Username,
		"email":    data.Email,
	}

	// Tạo token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret, err := s.GetJwtSecret()
	if err != nil {
		return nil, err
	}
	accessToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign JWT: %w", err)
	}

	refreshToken, err := (&RefreshTokenService{
		Size:          32,
		Cache:         s.Cache,
		TenantDb:      s.TenantDb,
		EncryptionKey: s.EncryptionKey,
		Context:       s.Context,
	}).GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	// Tạo OAuth2Token
	oauthToken := &OAuth2Token{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int64(tokenDuration.Seconds()),

		RefreshToken: refreshToken,
		Message:      "Token generated successfully",
		RoleId:       data.RoleId,
		UserId:       data.UserId,
		Username:     data.Username,
		Email:        data.Email,
	}

	return oauthToken, nil
}
