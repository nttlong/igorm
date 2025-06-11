package oauth

import (
	"context"
	"net/http"
	"strings"
	"unvs/internal/config"

	"dbx"
	"dynacall"

	"github.com/labstack/echo/v4"
)

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
	db.Open()
	return r
}

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
// @Success 200 {object} OAuth2TokenResponse "Đăng nhập thành công, trả về JWT token"
// @Failure 400 {object} ErrorResponse "Yêu cầu không hợp lệ (validation errors)"
// @Failure 401 {object} ErrorResponse "Thông tin đăng nhập không hợp lệ"
// @Failure 500 {object} ErrorResponse "Lỗi nội bộ server"
// @Router /oauth/token [post]
func (h *OAuthHandler) Token(c echo.Context) error {
	username := c.FormValue("username")
	if !strings.Contains(username, "@") { // Kiểm tra username có phải là email hay không
		return c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    "INVALID_USERNAME",
			Message: "Username must follow Tenant name",
		})
	}
	tanentName := strings.Split(username, "@")[1]
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

	ret, err := dynacall.Call(callPath, postData, struct {
		Tenant    string
		TenantDb  *dbx.DBXTenant
		Context   context.Context
		JwtSecret []byte
	}{
		Tenant:    "testDb",
		TenantDb:  createTenantDbx(tanentName),
		Context:   context.Background(),
		JwtSecret: config.GetJWTSecret(),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Internal server error",
		})
	}

	return c.JSON(http.StatusOK, ret)
}
