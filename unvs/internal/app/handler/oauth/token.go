package oauth

import (
	"context"
	"log"
	"net/http"
	"strings"
	"unvs/internal/config"

	"caching"
	"dbx"
	"dynacall"

	"github.com/labstack/echo/v4"
)

// LoginByFormSubmit xử lý yêu cầu HTTP POST để đăng nhập người dùng chỉ bằng form data.
// @Summary Đăng nhập người dùng bằng Form Submit (Username/Password)
// @Description Xác thực thông tin đăng nhập từ form data (username và mật khẩu) và trả về JWT token nếu thành công.
// @Tags Accounts
// @Accept x-www-form-urlencoded
// @Produce json
// @Param tenant formData string true "tenant code"
// @Param username formData string true "Tên người dùng (không phải email)"
// @Param password formData string true "Mật khẩu"
// @Param grant_type formData string false "Kiểu cấp quyền (thường là 'password' cho OAuth2, tùy chọn)"

// @Router /oauth/token [post]
func (h *OAuthHandler) Token(c echo.Context) error {
	username := c.FormValue("username")
	if !strings.Contains(username, "@") { // Kiểm tra username có phải là email hay không
		return c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    "INVALID_USERNAME",
			Message: "Username must follow Tenant name",
		})
	}
	tenantName := strings.Split(username, "@")[1]
	username = strings.Split(username, "@")[0]

	password := c.FormValue("password")
	grantType := c.FormValue("grant_type") // Lấy grant_type nếu có

	// Có thể thêm validation cho grant_type nếu nó là bắt buộc cho luồng OAuth2 cụ thể
	if grantType != "" && grantType != "password" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "UNSUPPORTED_GRANT_TYPE",
			Message: "grant_type 'password' supported only",
		})
	}
	callPath := "login@unvs.br.auth.users"
	postData := []interface{}{username, password}
	dbTenant, err := config.CreateTenantDbx(tenantName)
	if err != nil {
		log.Fatal(err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Internal server error",
		})
	}
	ret, err := dynacall.Call(callPath, postData, struct {
		Tenant   string
		TenantDb *dbx.DBXTenant
		Context  context.Context

		Cache         caching.Cache
		EncryptionKey string
	}{
		Tenant:        tenantName,
		EncryptionKey: config.AppConfigInstance.EncryptionKey,
		Cache:         config.GetCache(),
		Context:       c.Request().Context(),
		TenantDb:      dbTenant,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Internal server error",
		})
	}

	return c.JSON(http.StatusOK, ret)
}
