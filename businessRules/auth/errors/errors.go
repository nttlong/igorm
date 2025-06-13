package auth

import "fmt"

type AuthErrorCode int

const (
	ErrUnknown AuthErrorCode = iota
	ErrInvalidUsernameOrPassword
	ErrInvalidToken
	ErrInvalidRefreshToken
	ErrAccessDeny
)

func (e AuthErrorCode) String() string {
	switch e {
	case ErrUnknown:
		return "unknown error"
	case ErrInvalidUsernameOrPassword:
		return "invalid username or password"
	case ErrInvalidToken:
		return "invalid token"
	case ErrInvalidRefreshToken:
		return "invalid refresh token"
	default:
		return "unknown error"
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
