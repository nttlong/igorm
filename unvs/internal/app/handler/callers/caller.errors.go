package caller

import (
	"dbx"
	"dynacall"
	"net/http"
	"reflect"
	"strings"

	"github.com/labstack/echo/v4"
)

func (h *CallerHandler) CallHandlerErr(c echo.Context, err error, callerPath string) error {
	pkgPath := reflect.TypeOf(h).Elem().PkgPath() + "/Call"
	log := h.AppLogger.WithField("pkgPath", pkgPath).WithField("callerPath", callerPath)
	log.Error(err)
	if e, ok := err.(*dynacall.CallError); ok {
		if e.Code == dynacall.CallErrorCodeTokenExpired {
			return c.JSON(http.StatusUnauthorized, ErrorResponse{
				Code:    e.Code.String(),
				Message: e.Err.Error(),
			})
		}
		if e.Code == dynacall.CallErrorCodeAccessDenied {
			return c.JSON(http.StatusForbidden, ErrorResponse{
				Code:    e.Code.String(),
				Message: e.Err.Error(),
			})
		}
		if e.Code == dynacall.CallErrorCodeAuthenticationFailed {
			return c.JSON(http.StatusUnauthorized, ErrorResponse{
				Code:    e.Code.String(),
				Message: e.Err.Error(),
			})
		}
		if e.Code == dynacall.CallErrorCodeAuthenticationFailed {
			return c.JSON(http.StatusUnauthorized, ErrorResponse{
				Code:    e.Code.String(),
				Message: e.Err.Error(),
			})
		}

		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    e.Code.String(),
			Message: e.Err.Error(),
		})
	}
	if e, ok := err.(dynacall.CallError); ok {
		if strings.Contains(e.Err.Error(), "invalid args") {
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Code:    "INVALID_REQUEST_BODY",
				Message: "invalid request body",
			})
		}
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    e.Code.String(),
			Message: "Server error",
		})
	}
	if e, ok := err.(*dbx.DBXError); ok {
		if e.Code == dbx.DBXErrorCodeDuplicate {
			return c.JSON(http.StatusConflict, ErrorResponse{
				Code:    e.Code.String(),
				Fields:  e.Fields,
				Values:  e.Values,
				Message: e.Error(),
			})
		}
		if e.Code == dbx.DBXErrorCodeInvalidSize {

			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Code:    e.Code.String(),
				Message: e.Error(),
				MaxSize: e.MaxSize,
				Fields:  e.Fields,
			})
		}

		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    e.Code.String(),
			Message: e.Error(),
		})

	}
	return c.JSON(http.StatusInternalServerError, ErrorResponse{
		Code:    "SERVER_ERROR",
		Message: "Server error",
	})

}
