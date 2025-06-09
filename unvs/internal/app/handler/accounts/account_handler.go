// internal/app/handler/account/account_handler.go
package accounts // Tên package là 'account', theo tên thư mục

import (
	"dbx"

	_ "dbx" // Import model auth.User nếu cần trả về user info
	// Import errors để sử dụng errors.Is
	"unvs/internal/app/service/account" // Import account service
	"unvs/internal/model/auth"
	_ "unvs/internal/model/base"

	"github.com/golang-jwt/jwt/v4"
)

var x dbx.FullTextSearchColumn

// Request body struct for CreateAccount
type CreateAccountRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Response struct for successful account creation
type CreateAccountResponse struct {
	Message string     `json:"message"`
	User    *auth.User `json:"user,omitempty"`
}

// ErrorResponse struct for consistent error messages
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// AccountHandler là struct chứa dependency đến AccountService.
type AccountHandler struct {
	accountService *account.AccountService // Phụ thuộc vào AccountService
}

// NewAccountHandler tạo một instance mới của AccountHandler.
// AccountService được inject vào đây.
func NewAccountHandler(accSvc *account.AccountService) *AccountHandler {
	return &AccountHandler{accountService: accSvc}
}

// Response struct for successful Login
type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
	// Có thể bao gồm thêm thông tin người dùng nếu cần, nhưng không phải mật khẩu hash
	User *auth.User `json:"user,omitempty"`
}

// Request body struct for Login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserClaims defines the claims for the JWT token
type UserClaims struct {
	UserID   string `json:"userId"` // Sử dụng string cho UUID
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}
