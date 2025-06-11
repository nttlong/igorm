package caller

import (
	"context"
	"dbx"
	"dynacall"
	"net/http"

	"unvs/internal/config"
	_ "unvs/views"

	"github.com/labstack/echo/v4"
)

type CallerHandler struct {
}
type CallerRequest struct {
	Action   string        `json:"action"`
	Args     []interface{} `json:"args"`
	Tenant   string        `json:"tenant"`
	Language string        `json:"language"`
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

func createCfg() *dbx.Cfg {
	return &dbx.Cfg{
		Driver:   "mssql",
		Host:     "localhost",
		Port:     0,
		User:     "sa",
		Password: "123456",
		SSL:      false,
	}
}
func createDbx() *dbx.DBX {
	ret := dbx.NewDBX(*createCfg())
	return ret
}
func createTenantDbx(tenant string) *dbx.DBXTenant {
	db := createDbx()
	r, e := db.GetTenant(tenant)
	if e != nil {
		panic(e)
	}
	return r
}

// CallerHandler
// @summary CallerHandler
// @description CallerHandler
// @tags caller
// @accept json
// @produce json
// @Param request body CallerRequest true "CallerRequest"
// @router /callers/call [post]
// @Success 201 {object} CallerResponse "Response"
// @Security OAuth2Password
func (h *CallerHandler) Call(c echo.Context) error {

	req := new(CallerRequest)

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_REQUEST_BODY",
			Message: "Dữ liệu yêu cầu không hợp lệ",
		})
	}
	tenantDb, err := getTenantDb(req.Tenant)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_TENANT",
			Message: "Không tìm thấy tenant",
		})
	}
	retCall, err := dynacall.Call(req.Action, req.Args, struct {
		Tenant    string
		TenantDb  *dbx.DBXTenant
		Context   context.Context
		JwtSecret []byte
	}{
		Tenant:    req.Tenant,
		TenantDb:  tenantDb,
		Context:   c.Request().Context(),
		JwtSecret: config.GetJWTSecret(),
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
