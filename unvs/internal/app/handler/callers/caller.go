package caller

import (
	cache "caching"
	"context"
	"dbx"
	"dynacall"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"runtime/debug"
	"unvs/internal/config"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type CallerHandler struct {
	AppLogger *logrus.Logger
}
type CallerRequest struct {
	Args     interface{} `json:"args"`
	Tenant   string      `json:"tenant"`
	Language string      `json:"language"`
}
type ErrorResponse struct {
	Code    string   `json:"code"`
	Fields  []string `json:"fields"`
	Values  []string `json:"values"`
	Message string   `json:"message"`
}

// Response struct for successful account creation
type CallerResponse struct {
	Error   *ErrorResponse `json:"error,omitempty"`
	Results interface{}    `json:"results,omitempty"`
}

// CallerHandler
// @summary CallerHandler
// @description CallerHandler
// @tags caller
// @accept json
// @produce json
// @Param action path string true "The specific action to invoke (e.g., login, register, logout)"
// @Param request body CallerRequest true "CallerRequest"
// @router /invoke/{action} [post]
// @Success 201 {object} CallerResponse "Response"
// @Security OAuth2Password
func (h *CallerHandler) Call(c echo.Context) error {

	callerPath := c.Param("action")
	callerPath, err := url.PathUnescape(callerPath)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorResponse{
			Code:    "INVALID_REQUEST_URL",
			Message: "Invalid request URL",
		})
	}
	//req := new(CallerRequest)
	req, err := dynacall.NewRequestInstance(callerPath, reflect.TypeOf(CallerRequest{}))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Internal server error",
		})
	}

	if err := c.Bind(req.Data); err != nil {
		if eError, ok := err.(*echo.HTTPError); ok {
			h.AppLogger.Error(eError.Unwrap().Error())
			exampleData, errGet := dynacall.GetInputExampleCallerPath(callerPath)
			if errGet != nil {
				return c.JSON(http.StatusBadRequest, ErrorResponse{
					Code:    "INVALID_REQUEST_BODY",
					Message: "Invalid request body",
				})
			}
			if exampleData == nil || len(exampleData) == 0 {
				return c.JSON(http.StatusBadRequest, ErrorResponse{
					Code:    "INVALID_REQUEST_BODY",
					Message: "Invalid request body",
				})
			} else if len(exampleData) == 1 {
				jsonData, err := json.Marshal(exampleData[0])
				if err != nil {
					return c.JSON(http.StatusBadRequest, ErrorResponse{
						Code:    "INVALID_REQUEST_BODY",
						Message: "Invalid request body",
					})
				}
				msg := fmt.Sprintf("Invalid request body, Expected: %s", string(jsonData))
				return c.JSON(http.StatusBadRequest, ErrorResponse{
					Code:    "INVALID_REQUEST_BODY",
					Message: msg,
				})
			} else {
				jsonData, err := json.Marshal(exampleData)
				if err != nil {
					return c.JSON(http.StatusBadRequest, ErrorResponse{
						Code:    "INVALID_REQUEST_BODY",
						Message: "Invalid request body",
					})
				}
				msg := fmt.Sprintf("Invalid request body, Expected: %s", string(jsonData))
				return c.JSON(http.StatusBadRequest, ErrorResponse{
					Code:    "INVALID_REQUEST_BODY",
					Message: msg,
				})
			}

		}

		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_REQUEST_BODY",
			Message: "Invalid request body",
		})
	}

	tenantDb, err := config.CreateTenantDbx(req.GetString("Tenant"))
	if err != nil {
		log.Fatal(err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Internal server error",
		})
	}
	args := req.Get("Args")
	if omapDatak, ok := args.(interface{}); ok {

		defer func() {
			if r := recover(); r != nil {

				pkgPath := reflect.TypeOf(h).Elem().PkgPath() + "/Call"
				log := h.AppLogger.WithField("pkgPath", pkgPath).WithField("callerPath", callerPath)
				log.Errorf("Panic occurred: %v\n", r)
				log.Printf("Stack Trace:\n%s", debug.Stack())

			}
		}() // Gọi ngay lập tức hàm ẩn danh deferred
		retCall, err := dynacall.Call(callerPath, omapDatak, struct {
			Tenant        string
			TenantDb      *dbx.DBXTenant
			Context       context.Context
			EncryptionKey string

			Cache       cache.Cache
			AccessToken string
		}{
			Tenant:        req.GetString("Tenant"),
			TenantDb:      tenantDb,
			Context:       c.Request().Context(),
			EncryptionKey: config.AppConfigInstance.EncryptionKey,

			Cache:       config.GetCache(),
			AccessToken: c.Request().Header.Get("Authorization"),
		})
		if err != nil {
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
			if e, ok := err.(*dbx.DBXError); ok {
				if e.Code == dbx.DBXErrorCodeDuplicate {
					return c.JSON(http.StatusConflict, ErrorResponse{
						Code:    e.Code.String(),
						Fields:  e.Fields,
						Values:  e.Values,
						Message: e.Error(),
					})
				}

			}
			return c.JSON(http.StatusInternalServerError, ErrorResponse{
				Code:    "SERVER_ERROR",
				Message: "Server error",
			})
		}
		// err = retCall[1].(error)
		// if err != nil {
		// 	return c.JSON(http.StatusBadRequest, err)
		// }
		return c.JSON(http.StatusOK, retCall)
	}
	return c.JSON(http.StatusBadRequest, ErrorResponse{
		Code:    "INVALID_REQUEST_BODY",
		Message: "Dữ liệu yêu cầu không hợp lệ",
	})

}
