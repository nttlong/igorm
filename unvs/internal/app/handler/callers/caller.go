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
	Args interface{} `json:"args"`
	// Tenant   string      `json:"tenant"`
	// Language string      `json:"language"`
}
type ErrorResponse struct {
	Code    string   `json:"code"`
	Fields  []string `json:"fields"`
	Values  []string `json:"values"`
	Message string   `json:"message"`
	MaxSize int      `json:"maxSize"`
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
// @Param feature query string true "The specific id of feature. Each UI at frontend will have a unique feature id and must be approve by backend team."
// @Param action query string true "The specific action to invoke (e.g., login, register, logout)"
// @Param module query string true "The specific module to invoke (e.g., unvs.br.auth.roles, unvs.br.auth.uusers, ...)"
// @Param tenant query string true "The specific tenant to invoke (e.g., default, name, ...)"
// @Param lan query string true "The specific language to invoke (e.g., en, pt, ...)"
// @Param request body CallerRequest true "CallerRequest"
// @router /invoke [post]
// @Success 201 {object} CallerResponse "Response"
// @Security OAuth2Password
func (h *CallerHandler) Call(c echo.Context) error {
	info, err := ExtractRequireQueryStrings(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	fmt.Println(info)
	callerPath := info.Action + "@" + info.Module
	req, err := dynacall.NewInvoker(callerPath)
	err = req.New(callerPath)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	// req, err := dynacall.NewRequestInstance(callerPath)
	// args := req.Get("Args")
	// tpy := reflect.TypeOf(args)
	// for i := 0; i < tpy.NumField(); i++ {
	// 	field := tpy.Field(i)
	// 	fmt.Println(field.Name)
	// }
	// fmt.Println(reflect.TypeOf(args).Name())

	// jsonBff, err := json.Marshal(req.Data)
	// fmt.Print(string(jsonBff))

	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_REQUEST_BODY",
			Message: "Invalid request body",
		})
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Internal server error",
		})
	}
	// vv := &struct { // Định nghĩa struct ẩn danh ngay đây
	// 	Username string
	// 	Password string
	// }{
	// 	Username: "admin",  // Gán giá trị cụ thể
	// 	Password: "123456", // Gán giá trị cụ thể
	// }

	if err := c.Bind(req.Args); err != nil {

		if eError, ok := err.(*echo.HTTPError); ok {
			h.AppLogger.Error(eError.Unwrap().Error())
			exampleData, errGet := dynacall.GetInputExampleCallerPath(callerPath)
			if errGet != nil {
				return c.JSON(http.StatusBadRequest, ErrorResponse{
					Code:    "INVALID_REQUEST_BODY",
					Message: "Invalid request body",
				})
			}
			if len(exampleData) == 0 {
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

	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_REQUEST_BODY",
			Message: "Invalid request body",
		})
	}

	tenantDb, err := config.CreateTenantDbx(info.Tenant)
	if err != nil {
		appLogger := h.AppLogger.WithField("pkgPath", reflect.TypeOf(h).Elem().PkgPath()+"/Call")
		appLogger.Error(err)
		return c.JSON(http.StatusForbidden, ErrorResponse{
			Code:    "FORBIDDEN",
			Message: "Forbidden",
		})
	}
	err = tenantDb.Open()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "internal server error",
		})

	}

	if err != nil {
		log.Fatal(err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Internal server error",
		})
	}

	// args := req.Get("Args")
	defer func() {
		if r := recover(); r != nil {

			pkgPath := reflect.TypeOf(h).Elem().PkgPath() + "/Call"
			log := h.AppLogger.WithField("pkgPath", pkgPath).WithField("callerPath", callerPath)
			log.Errorf("Panic occurred: %v\n", r)
			log.Printf("Stack Trace:\n%s", debug.Stack())

		}
	}() // Gọi ngay lập tức hàm ẩn danh deferred
	injectorCaller := dynacall.NewCaller(callerPath, &struct {
		Tenant        string
		TenantDb      *dbx.DBXTenant
		Context       context.Context
		EncryptionKey string
		Language      string

		Cache       cache.Cache
		AccessToken string
		FeatureId   string
	}{
		Tenant:        info.Tenant,
		TenantDb:      tenantDb,
		Context:       c.Request().Context(),
		EncryptionKey: config.AppConfigInstance.EncryptionKey,
		Language:      info.Lan,

		Cache:       config.GetCache(),
		AccessToken: c.Request().Header.Get("Authorization"),
		FeatureId:   info.Feature,
	})
	retCall, err := injectorCaller(req.Args)
	if err != nil {
		return h.CallHandlerErr(c, err, callerPath)
	}
	return c.JSON(http.StatusOK, CallerResponse{
		Results: retCall,
	})

}
