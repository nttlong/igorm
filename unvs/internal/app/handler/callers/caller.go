package caller

import (
	"context"
	"dbx"
	"dynacall"
	"log"
	"net/http"
	"net/url"
	"reflect"

	cache "caching"
	"unvs/internal/config"

	"github.com/labstack/echo/v4"
)

type CallerHandler struct {
}
type CallerRequest struct {
	Args     interface{} `json:"args"`
	Tenant   string      `json:"tenant"`
	Language string      `json:"language"`
}
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
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

	if err := c.Bind(req.Data); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_REQUEST_BODY",
			Message: "Dữ liệu yêu cầu không hợp lệ",
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

		retCall, err := dynacall.Call(callerPath, omapDatak, struct {
			Tenant    string
			TenantDb  *dbx.DBXTenant
			Context   context.Context
			JwtSecret []byte
			Cache     cache.Cache
		}{
			Tenant:    req.GetString("Tenant"),
			TenantDb:  tenantDb,
			Context:   c.Request().Context(),
			JwtSecret: config.GetJWTSecret(),
			Cache:     config.GetCache(),
		})
		if err != nil {
			if e, ok := err.(dynacall.CallError); ok {
				return c.JSON(http.StatusBadRequest, ErrorResponse{
					Code:    e.Code.String(),
					Message: e.Err.Error(),
				})
			}
			return c.JSON(http.StatusBadRequest, err)
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
