package caller

import (
	cache "caching"
	"context"
	"dbx"
	"dynacall"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"reflect"
	"runtime/debug"
	"unvs/internal/config"

	"github.com/labstack/echo/v4"
)

// FormRequest đại diện cho dữ liệu form và các file upload từ người dùng.
// Các trường form text có thể được bind tự động.
// Các trường file upload sẽ được xử lý riêng và lưu tham chiếu vào đây.
type FormRequest struct {
	Data string `json:"data" form:"data"`

	// Các trường file upload (không thể bind trực tiếp bằng tag 'form' hay 'json' từ request body)
	// Bạn sẽ phải lấy file từ c.FormFile hoặc c.MultipartForm và gán vào đây sau.
	UploadedFiles []*multipart.FileHeader `json:"-" form:"-"` // Dấu "-" để bỏ qua binding tự động
}

// CallerResponse (Giả định struct phản hồi của bạn)
type FormResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Submit handler
// @summary Submit handler for form data and file uploads
// @description Handles form data submission including file uploads.
// @tags caller
// @accept multipart/form-data
// @produce json
// @Param feature query string true "The specific id of feature. Each UI at frontend will have a unique feature id and must be approve by backend team."
// @Param tenant query string true "The specific tenant to invoke (e.g., default, name, ...)"
// @Param module query string true "The specific module to invoke (e.g., unvs.br.auth.roles, unvs.br.auth.uusers, ...)"
// @Param action query string true "The specific action to invoke (e.g., login, register, logout)"
// @Param lan query string true "The specific language to invoke (e.g., en, pt, ...)"
// @Param data formData string true "JSON stringify from browser" default({"code":"R003","name":"test3","description":"example description"})
// @Param files formData file false "One or more files to upload"
// @router /invoke-form [post]
// @Success 201 {object} CallerResponse "Response"
// @Security OAuth2Password
func (h *CallerHandler) FormSubmit(c echo.Context) error {
	// get data from request body
	info, err := ExtractRequireQueryStrings(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	fmt.Println(info)
	callerPath := info.Action + "@" + info.Module
	req, err := dynacall.NewRequestInstance(callerPath, reflect.TypeOf(CallerRequest{}))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	jsonData, _ := json.Marshal(req)

	fmt.Println(jsonData)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	data := "{\"Args\":" + c.FormValue("data") + "}"

	err = json.Unmarshal([]byte(data), &req.Data)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
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
	if omapDatak, ok := req.Data.(interface{}); ok {

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
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		response := FormResponse{
			Status:  "success",
			Message: "Call success",
			Data:    retCall,
		}
		return c.JSON(http.StatusOK, response)
	}
	return echo.NewHTTPError(http.StatusBadRequest, "Invalid data format")
}
