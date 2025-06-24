package errors

import "fmt"

type ErrorCode int

const (
	Error_Unknown ErrorCode = iota

	Error_AccessDenied
	Error_TokenExpired
	Error_LoginFailed
)

func (e ErrorCode) String() string {
	switch e {
	case Error_Unknown:
		return "UNKNOWN"
	case Error_AccessDenied:
		return "ACCESS_DENIED"
	case Error_TokenExpired:
		return "TOKEN_EXPIRED"
	case Error_LoginFailed:
		return "LOGIN_FAILED"
	default:
		return "UNKNOWN"
	}
}

type CoreError struct {
	Code    ErrorCode
	Message string
	Err     error
}

func (e *CoreError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}
