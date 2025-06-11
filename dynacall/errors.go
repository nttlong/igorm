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
)

func (e CallErrorCode) String() string {
	switch e {
	case CallErrorCodeUnknown:
		return "unknown error"
	case CallErrorCodeInvalidParams:
		return "invalid parameters"
	case CallErrorCodeInvalidReturn:
		return "invalid return value"
	case CallErrorCodeTimeout:
		return "timeout"
	case CallErrorCodeInternalError:
		return "internal error"
	case CallErrorCodeInvalidCallerPath:
		return "invalid caller path"
	case CallErrorCodeCallerPathNotFound:
		return "caller path not found"
	case CallerSystemError:
		return "caller system error"
	default:
		return "unknown error"
	}
}

type CallError struct {
	Err  error
	Code CallErrorCode
}

func (e CallError) Error() string {
	return fmt.Sprintf("dynacall error: %s (%s)", e.Err.Error(), e.Code.String())
}
