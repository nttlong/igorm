package accounts

// import (
// 	"net/http"
// 	"strings"

// 	"github.com/labstack/echo/v4"
// )

// // LoginByFormSubmit xử lý yêu cầu HTTP POST để đăng nhập người dùng chỉ bằng form data.
// // @Summary Đăng nhập người dùng bằng Form Submit (Username/Password)
// // @Description Xác thực thông tin đăng nhập từ form data (username và mật khẩu) và trả về JWT token nếu thành công.
// // @Tags Accounts
// // @Accept x-www-form-urlencoded
// // @Produce json
// // @Param username formData string true "Tên người dùng (không phải email)" // <<< Sử dụng formData cho Swagger
// // @Param password formData string true "Mật khẩu"
// // @Param grant_type formData string false "Kiểu cấp quyền (thường là 'password' cho OAuth2, tùy chọn)"
// // @Success 200 {object} LoginResponse "Đăng nhập thành công, trả về JWT token"
// // @Failure 400 {object} ErrorResponse "Yêu cầu không hợp lệ (validation errors)"
// // @Failure 401 {object} ErrorResponse "Thông tin đăng nhập không hợp lệ"
// // @Failure 500 {object} ErrorResponse "Lỗi nội bộ server"
// // @Router /oauth/token [post]
// func (h *AccountHandler) LoginByFormSubmit(c echo.Context) error {
// 	username := c.FormValue("username")
// 	password := c.FormValue("password")
// 	grantType := c.FormValue("grant_type") // Lấy grant_type nếu có

// 	// Có thể thêm validation cho grant_type nếu nó là bắt buộc cho luồng OAuth2 cụ thể
// 	if grantType != "" && grantType != "password" {
// 		return c.JSON(http.StatusBadRequest, ErrorResponse{
// 			Code:    "UNSUPPORTED_GRANT_TYPE",
// 			Message: "Chỉ hỗ trợ grant_type 'password'",
// 		})
// 	}

// 	// Gọi helper method để xử lý logic đăng nhập chính
// 	response, err := h.performLoginRequest(c.Request().Context(), username, password)
// 	if err != nil {
// 		// Phân tích lỗi từ performLoginRequest
// 		if strings.HasPrefix(err.Error(), "MISSING_CREDENTIALS") {
// 			return c.JSON(http.StatusBadRequest, ErrorResponse{
// 				Code:    "MISSING_CREDENTIALS",
// 				Message: strings.TrimPrefix(err.Error(), "MISSING_CREDENTIALS: "),
// 			})
// 		} else if strings.HasPrefix(err.Error(), "INVALID_CREDENTIALS") {
// 			return c.JSON(http.StatusUnauthorized, ErrorResponse{
// 				Code:    "INVALID_CREDENTIALS",
// 				Message: strings.TrimPrefix(err.Error(), "INVALID_CREDENTIALS: "),
// 			})
// 		} else if strings.HasPrefix(err.Error(), "TOKEN_GENERATION_FAILED") {
// 			return c.JSON(http.StatusInternalServerError, ErrorResponse{
// 				Code:    "TOKEN_GENERATION_FAILED",
// 				Message: strings.TrimPrefix(err.Error(), "TOKEN_GENERATION_FAILED: "),
// 			})
// 		} else { // AUTHENTICATION_ERROR hoặc lỗi không xác định khác
// 			return c.JSON(http.StatusInternalServerError, ErrorResponse{
// 				Code:    "AUTHENTICATION_ERROR",
// 				Message: "Lỗi xác thực không xác định",
// 			})
// 		}
// 	}

// 	return c.JSON(http.StatusOK, response)
// }
