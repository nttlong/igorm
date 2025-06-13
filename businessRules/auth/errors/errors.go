package auth

import "fmt"

type AuthErrorCode int

const (
	ErrUnknown AuthErrorCode = iota
	ErrInvalidUsernameOrPassword
	ErrInvalidToken
	ErrInvalidRefreshToken
	ErrAccessDeny
	ErrTokenExpired
)

func (e AuthErrorCode) String() string {
	switch e {
	case ErrUnknown:
		return "UNKNOWN_ERROR"
	case ErrInvalidUsernameOrPassword:
		return "INVALID_USERNAME_OR_PASSWORD"
	case ErrInvalidToken:
		return "INVALID_TOKEN"
	case ErrInvalidRefreshToken:
		return "INVALID_REFRESH_TOKEN"
	default:
		return "UNKNOWN_ERROR"
	}
}

type AuthError struct {
	Code    AuthErrorCode
	Message string
	Err     error
}

func (e *AuthError) Error() string {

	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}
