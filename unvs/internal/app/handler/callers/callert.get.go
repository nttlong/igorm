package caller

import (
	"caching"
	"context"
	"dbx"
	"dynacall"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"runtime/debug"
	"strings"
	"unvs/internal/config"

	"github.com/labstack/echo/v4"
)

// CallGet
// @summary Lấy thông tin dựa trên tenant, module, action, feature, ngôn ngữ và các phân đoạn đường dẫn bổ sung.
// @description API này cho phép gọi các tính năng cụ thể cho từng tenant và ngôn ngữ, với khả năng mở rộng đường dẫn.
// @tags caller
// @accept json
// @produce json
// @Param tenant path string true "The specific tenant to invoke (e.g., default, name, ...)"
// @Param module path string true "The module name (e.g., users, products, auth)"
// @Param action path string true "The action name (e.g., list, create, detail, login)"
// @Param optionalPath path string false "Optional additional path segments (e.g., 'sub/item/id'). This parameter captures all remaining path segments."
// @Param feature query string true "The specific ID of the feature. Each UI at frontend will have a unique feature ID and must be approved by the backend team."
// @Param lan query string true "The specific language to invoke (e.g., en, vi, pt, ...)"
// @router /get/{tenant}/{module}/{action}/{optionalPath} [get]
// @Success 200 {object} CallerResponse "Successful response with requested parameters"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
func (h *CallerHandler) CallGet(c echo.Context) error {
	tenant := c.Param("tenant")
	module := c.Param("module")
	action := c.Param("action")
	params := c.Param("*")
	path := c.Request().URL.String()
	fmt.Println(path)

	// Lấy tham số từ query string
	feature := c.Param("feature")
	lan := c.Param("lan")
	info := ExtractInfo{
		Feature: feature,
		Tenant:  tenant,
		Lan:     lan,
		Action:  action,
		Module:  module,
	}

	fmt.Println(info)
	callerPath := info.Action + "@" + info.Module
	req, err := dynacall.NewInvoker(callerPath)
	err = req.New(callerPath)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	req.Args = strings.Split(params, "/")
	req.Args = []interface{}{"vi"}
	fmt.Println(params)

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
	callPath := info.Action + "@" + info.Module
	// args := req.Get("Args")
	defer func() {
		if r := recover(); r != nil {

			pkgPath := reflect.TypeOf(h).Elem().PkgPath() + "/Call"
			log := h.AppLogger.WithField("pkgPath", pkgPath).WithField("callerPath", callerPath)
			log.Errorf("Panic occurred: %v\n", r)
			log.Printf("Stack Trace:\n%s", debug.Stack())

		}
	}() // Gọi ngay lập tức hàm ẩn danh deferred
	postData := []interface{}{"vi"}

	ret, err := dynacall.Call(callPath, postData, struct {
		Tenant   string
		TenantDb *dbx.DBXTenant
		Context  context.Context

		Cache         caching.Cache
		EncryptionKey string
	}{
		Tenant:        tenant,
		EncryptionKey: config.AppConfigInstance.EncryptionKey,
		Cache:         config.GetCache(),
		Context:       c.Request().Context(),
		TenantDb:      tenantDb,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Internal server error",
		})
	}

	return c.JSON(http.StatusOK, ret)

}
