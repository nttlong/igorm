package account

import (
	"errors"
	"fmt"
	"net/http"
	"unvs/internal/app/service/account"

	"github.com/labstack/echo/v4"
)

// CreateAccount xử lý yêu cầu HTTP POST để tạo tài khoản mới.
// @Summary Tạo một tài khoản người dùng mới
// @Description Tạo một tài khoản người dùng mới với username, email và mật khẩu.
// @Tags Accounts
// @Accept json
// @Produce json
// @Param request body CreateAccountRequest true "Thông tin tài khoản cần tạo"
// @Success 201 {object} CreateAccountResponse "Tạo tài khoản thành công"
// @Failure 400 {object} ErrorResponse "Yêu cầu không hợp lệ (validation errors)"
// @Failure 409 {object} ErrorResponse "Email đã tồn tại"
// @Failure 500 {object} ErrorResponse "Lỗi nội bộ server"
// @Router /accounts/create [post]
// @Security OAuth2Password
func (h *AccountHandler) CreateAccount(c echo.Context) error {
	req := new(CreateAccountRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_REQUEST_BODY",
			Message: "Dữ liệu yêu cầu không hợp lệ",
		})
	}

	// Gọi phương thức nghiệp vụ từ Service Layer
	newUser, err := h.accountService.CreateAccount(c.Request().Context(), req.Username, req.Email, req.Password)
	if err != nil {
		// Xử lý các lỗi nghiệp vụ từ Service Layer
		switch {
		case errors.Is(err, account.ErrEmailEmpty):
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Code:    "EMAIL_EMPTY",
				Message: "Email không được để trống",
			})
		case errors.Is(err, account.ErrPasswordEmpty):
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Code:    "PASSWORD_EMPTY",
				Message: "Mật khẩu không được để trống",
			})
		case errors.Is(err, account.ErrPasswordTooShort):
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Code:    "PASSWORD_TOO_SHORT",
				Message: "Mật khẩu phải có ít nhất 6 ký tự",
			})
		case errors.Is(err, account.ErrEmailAlreadyUsed):
			return c.JSON(http.StatusConflict, ErrorResponse{ // 409 Conflict
				Code:    "EMAIL_ALREADY_USED",
				Message: "Email đã được sử dụng",
			})
		case errors.Is(err, account.ErrHashPassword):
			// Đây là lỗi nội bộ, không nên trả về chi tiết cho client
			return c.JSON(http.StatusInternalServerError, ErrorResponse{
				Code:    "INTERNAL_ERROR",
				Message: "Đã xảy ra lỗi nội bộ khi xử lý mật khẩu",
			})
		case errors.Is(err, account.ErrCreateAccountFail):
			// Đây là lỗi tổng quát khi tạo tài khoản, có thể do lỗi DB khác không phải duplicate
			return c.JSON(http.StatusInternalServerError, ErrorResponse{
				Code:    "CREATE_ACCOUNT_FAILED",
				Message: "Không thể tạo tài khoản, vui lòng thử lại sau",
			})
		case errors.Is(err, account.ErrUsernameAlreadyUsed):
			// Đây là lỗi tổng quát khi tạo tài khoản, có thể do lỗi DB khác không phải duplicate
			return c.JSON(http.StatusConflict, ErrorResponse{
				Code:    "USERNAME_ALREADY_USED",
				Message: "Tên tài khoản đã được sử dụng, vui lòng chọn tên khác",
			})
		default:
			// Xử lý các lỗi không xác định
			return c.JSON(http.StatusInternalServerError, ErrorResponse{
				Code:    "UNKNOWN_ERROR",
				Message: fmt.Sprintf("Đã xảy ra lỗi không xác định: %v", err),
			})
		}
	}

	// Trả về phản hồi thành công
	return c.JSON(http.StatusCreated, CreateAccountResponse{
		Message: "Tài khoản đã được tạo thành công",
		User:    newUser, // Trả về thông tin người dùng đã tạo (trừ password hash)
	})
}
