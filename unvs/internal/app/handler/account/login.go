package account

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// Login xử lý yêu cầu HTTP POST để đăng nhập người dùng (JSON hoặc Form).
// @Summary Đăng nhập người dùng và nhận JWT token (JSON/Form)
// @Description Xác thực thông tin đăng nhập của người dùng (email hoặc username) và trả về JWT token nếu thành công.
// @Tags Accounts
// @Accept json
// @Accept x-www-form-urlencoded
// @Produce json
// @Param request body LoginRequest true "Thông tin đăng nhập (email/username và mật khẩu)"
// @Success 200 {object} LoginResponse "Đăng nhập thành công, trả về JWT token"
// @Failure 400 {object} ErrorResponse "Yêu cầu không hợp lệ (validation errors)"
// @Failure 401 {object} ErrorResponse "Thông tin đăng nhập không hợp lệ"
// @Failure 500 {object} ErrorResponse "Lỗi nội bộ server"
// @Router /accounts/login [post]
func (h *AccountHandler) Login(c echo.Context) error {
	req := new(LoginRequest)
	// Echo's c.Bind() with `json` and `form` tags can handle both JSON and x-www-form-urlencoded
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_REQUEST_BODY",
			Message: "Dữ liệu yêu cầu không hợp lệ",
		})
	}

	// Gọi helper method để xử lý logic đăng nhập chính
	response, err := h.performLoginRequest(c.Request().Context(), req.Username, req.Password)
	if err != nil {
		// Phân tích lỗi từ performLoginRequest
		if strings.HasPrefix(err.Error(), "MISSING_CREDENTIALS") {
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Code:    "MISSING_CREDENTIALS",
				Message: strings.TrimPrefix(err.Error(), "MISSING_CREDENTIALS: "),
			})
		} else if strings.HasPrefix(err.Error(), "INVALID_CREDENTIALS") {
			return c.JSON(http.StatusUnauthorized, ErrorResponse{
				Code:    "INVALID_CREDENTIALS",
				Message: strings.TrimPrefix(err.Error(), "INVALID_CREDENTIALS: "),
			})
		} else if strings.HasPrefix(err.Error(), "TOKEN_GENERATION_FAILED") {
			return c.JSON(http.StatusInternalServerError, ErrorResponse{
				Code:    "TOKEN_GENERATION_FAILED",
				Message: strings.TrimPrefix(err.Error(), "TOKEN_GENERATION_FAILED: "),
			})
		} else { // AUTHENTICATION_ERROR hoặc lỗi không xác định khác
			return c.JSON(http.StatusInternalServerError, ErrorResponse{
				Code:    "AUTHENTICATION_ERROR",
				Message: "Lỗi xác thực không xác định",
			})
		}
	}

	return c.JSON(http.StatusOK, response)
}
