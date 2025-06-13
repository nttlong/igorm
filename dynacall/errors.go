/*
 this file contains error definitions for dynacall package

*/

package dynacall

import "fmt"

type CallErrorCode int

const (
	CallErrorCodeUnknown CallErrorCode = iota
	CallErrorCodeInvalidParams
	CallErrorCodeInvalidReturn
	CallErrorCodeTimeout
	CallErrorCodeInternalError
	CallErrorCodeInvalidCallerPath
	CallErrorCodeCallerPathNotFound
	CallerSystemError
	CallErrorCodeInvalidArgs
	CallErrorCodeTokenExpired
	CallErrorCodeTokenInvalid
	CallErrorCodeTokenNotFound
	CallErrorCodeTokenInvalidFormat
	CallErrorCodeTokenInvalidIssuer
	CallErrorCodeTokenInvalidAudience
	CallErrorCodeTokenInvalidSubject
	CallErrorCodeTokenInvalidSignature
	CallErrorCodeTokenInvalidExpiry
	CallErrorCodeAuthorizationFailed
	CallErrorCodeAuthenticationFailed
	CallErrorCodeAccessDenied
)

func (e CallErrorCode) String() string {
	switch e {
	case CallErrorCodeUnknown:
		return "UNKNOWN"
	case CallErrorCodeInvalidParams:
		return "INVALID_PARAMS"
	case CallErrorCodeInvalidReturn:
		return "INVALID_RETURN"
	case CallErrorCodeTimeout:
		return "TIMEOUT"
	case CallErrorCodeInternalError:
		return "INTERNAL_ERROR"
	case CallErrorCodeInvalidCallerPath:
		return "INVALID_CALLER_PATH"
	case CallErrorCodeCallerPathNotFound:
		return "CALLER_PATH_NOT_FOUND"
	case CallerSystemError:
		return "CALLER_SYSTEM_ERROR"
	case CallErrorCodeInvalidArgs:
		return "INVALID_ARGS"
	case CallErrorCodeTokenExpired:
		return "TOKEN_EXPIRED"
	case CallErrorCodeTokenInvalid:
		return "TOKEN_INVALID"
	case CallErrorCodeTokenNotFound:
		return "TOKEN_NOT_FOUND"
	case CallErrorCodeTokenInvalidFormat:
		return "TOKEN_INVALID_FORMAT"
	case CallErrorCodeTokenInvalidIssuer:
		return "TOKEN_INVALID_ISSUER"
	case CallErrorCodeTokenInvalidAudience:
		return "TOKEN_INVALID_AUDIENCE"
	case CallErrorCodeTokenInvalidSubject:
		return "TOKEN_INVALID_SUBJECT"
	case CallErrorCodeTokenInvalidSignature:
		return "TOKEN_INVALID_SIGNATURE"
	case CallErrorCodeTokenInvalidExpiry:
		return "TOKEN_INVALID_EXPIRY"
	case CallErrorCodeAuthorizationFailed:
		return "AUTHORIZATION_FAILED"
	case CallErrorCodeAuthenticationFailed:
		return "AUTHENTICATION_FAILED"
	case CallErrorCodeAccessDenied:
		return "ACCESS_DENIED"
	default:
		return "UNKNOWN"
	}
}

type CallError struct {
	Err  error
	Code CallErrorCode
}

func (e CallError) Error() string {
	return fmt.Sprintf("dynacall error: %s (%s)", e.Err.Error(), e.Code.String())
}
