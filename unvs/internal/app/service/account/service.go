// internal/app/service/account/service.go
package account // Tên package là 'account', theo tên thư mục

import (
	"context"
	"errors"
	"fmt"
	"strings" // Cần import strings để sử dụng strings.Contains
	"time"

	"dbx"                                        // Import package dbx để kiểm tra lỗi DBXError
	userRepo "unvs/internal/app/repository/user" // Import user repository
	"unvs/internal/model/auth"                   // Import struct User từ model/auth

	// Import BaseModel từ model/base
	"github.com/google/uuid"     // Import thư viện uuid
	"golang.org/x/crypto/bcrypt" // Import thư viện bcrypt để so sánh mật khẩu
)

// === Các biến lỗi tùy chỉnh cho tầng Service (Domain-Specific Errors) ===
// Các lỗi này sử dụng mã (code) để dễ dàng mapping với các thông báo đa ngôn ngữ ở tầng Handler/Presentation.
var (
	ErrEmailEmpty          = errors.New("EMAIL_EMPTY")
	ErrPasswordEmpty       = errors.New("PASSWORD_EMPTY")
	ErrPasswordTooShort    = errors.New("PASSWORD_TOO_SHORT")
	ErrHashPassword        = errors.New("HASH_PASSWORD_FAILED")
	ErrEmailAlreadyUsed    = errors.New("EMAIL_ALREADY_USED")
	ErrCreateAccountFail   = errors.New("CREATE_ACCOUNT_FAILED")
	ErrInvalidCredentials  = errors.New("INVALID_CREDENTIALS")
	ErrUserNotFound        = errors.New("USER_NOT_FOUND")
	ErrUpdateUserFail      = errors.New("UPDATE_USER_FAILED")
	ErrDeleteUserFail      = errors.New("DELETE_USER_FAILED")
	ErrUsernameAlreadyUsed = errors.New("USERNAME_ALREADY_USED")
)

// AccountService là struct chứa logic nghiệp vụ cho các tác vụ liên quan đến tài khoản.
// Nó phụ thuộc vào UserRepository interface để truy cập dữ liệu.
type AccountService struct {
	userRepo userRepo.UserRepository // Phụ thuộc vào userRepo.UserRepository interface
}

// NewAccountService tạo một instance mới của AccountService.
// UserRepository được inject vào đây.
func NewAccountService(repo userRepo.UserRepository) *AccountService {
	return &AccountService{userRepo: repo}
}

// CreateAccount là phương thức nghiệp vụ để tạo một tài khoản mới.
// Nó xử lý logic validation và gọi repository để lưu trữ dữ liệu.
func (s *AccountService) CreateAccount(ctx context.Context, username string, email, password string) (*auth.User, error) {
	// 1. Logic nghiệp vụ: Validate input
	if email == "" {
		return nil, ErrEmailEmpty
	}
	if password == "" {
		return nil, ErrPasswordEmpty
	}
	if len(password) < 6 { // Ví dụ: mật khẩu tối thiểu 6 ký tự
		return nil, ErrPasswordTooShort
	}

	// 2. Băm mật khẩu.
	// Đây là nơi thực hiện logic băm mật khẩu theo nghiệp vụ của ứng dụng.
	// Repository (user_repo.go) KHÔNG NÊN băm mật khẩu.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrHashPassword, err) // Wrap lỗi gốc
	}

	// 3. Tạo đối tượng User mới
	newUser := &auth.User{
		UserId:       uuid.New().String(), // ĐÃ CHUYỂN LOGIC TẠO UUID VÀO ĐÂY (Service Layer)
		Username:     username,            // Gán username từ tham số
		Email:        email,
		PasswordHash: string(hashedPassword), // Gán mật khẩu đã băm
		CreatedAt:    time.Now(),
		CreatedBy:    "system",
	}

	// 4. Lưu vào DB thông qua Repository
	// Repository sẽ nhận newUser với PasswordHash đã được băm và UserId đã được gán.
	// ĐẢM BẢO user_repo.go KHÔNG CÓ logic băm mật khẩu và KHÔNG GÁN UUID.
	if err := s.userRepo.CreateUser(ctx, newUser); err != nil {
		// --- THÊM DÒNG NÀY ĐỂ GỠ LỖI: In ra chi tiết lỗi từ Repository ---
		fmt.Printf("DEBUG: Error from userRepo.CreateUser: Type=%T, Value=%v\n", err, err)

		var dbxErr *dbx.DBXError     // Khai báo biến kiểu DBXError
		if errors.As(err, &dbxErr) { // Sử dụng errors.As để kiểm tra lỗi có phải là DBXError không
			// --- THÊM DÒNG NÀY ĐỂ GỠ LỖI: In ra chi tiết DBXError ---
			fmt.Printf("DEBUG: Error is *dbx.DBXError. Code=%s, TableName=%s, ConstraintName=%s, Fields=%v, Values=%v, Message=%s\n",
				dbxErr.Code.String(), dbxErr.TableName, dbxErr.ConstraintName, dbxErr.Fields, dbxErr.Values, dbxErr.Message)

			if dbxErr.Code == dbx.DBXErrorCodeDuplicate {
				// Nếu là lỗi trùng lặp, kiểm tra trường Email
				for _, field := range dbxErr.Fields {
					// So sánh field với "Email" (hoặc "email" nếu tên cột DB là chữ thường)
					if strings.EqualFold(field, "Email") { // Sử dụng EqualFold để không phân biệt chữ hoa/thường
						return nil, ErrEmailAlreadyUsed // Trả về lỗi nghiệp vụ rõ ràng
					}
					if strings.EqualFold(field, "Username") { // Sử dụng EqualFold
						return nil, ErrUsernameAlreadyUsed // Trả về lỗi nghiệp vụ rõ ràng
					}
				}
				// Xử lý các trường hợp trùng lặp khác nếu có, hoặc trả về lỗi tổng quát
				return nil, fmt.Errorf("%w: %s", ErrCreateAccountFail, dbxErr.Message) // Wrap lỗi DBX
			}
			// Xử lý các loại lỗi DBX khác nếu cần, ví dụ:
			if dbxErr.Code == dbx.DBXErrorCodeInvalidSize {
				return nil, fmt.Errorf("%w: invalid size: %s", ErrCreateAccountFail, dbxErr.Message)
			}
			return nil, fmt.Errorf("%w: dbx error: %s", ErrCreateAccountFail, dbxErr.Message) // Wrap lỗi DBX khác
		}
		// Xử lý các lỗi khác không phải DBXError từ repository
		return nil, fmt.Errorf("%w: %v", ErrCreateAccountFail, err) // Wrap lỗi gốc
	}

	return newUser, nil
}

// AuthenticateUser là phương thức nghiệp vụ để xác thực người dùng.
func (s *AccountService) AuthenticateUser(ctx context.Context, username, password string) (*auth.User, error) {
	// 1. Lấy người dùng từ DB bằng email
	user, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		// Kiểm tra lỗi "not found" từ repository (dựa vào thông báo lỗi từ repo)
		if strings.Contains(err.Error(), "not found") {
			return nil, ErrInvalidCredentials // Dịch thành lỗi nghiệp vụ
		}
		return nil, fmt.Errorf("%w: failed to get user: %v", ErrInvalidCredentials, err) // Wrap lỗi gốc
	}

	// 2. So sánh mật khẩu băm với mật khẩu thô
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		// bcrypt.CompareHashAndPassword trả về lỗi nếu mật khẩu không khớp
		return nil, ErrInvalidCredentials // Dịch thành lỗi nghiệp vụ
	}

	return user, nil
}

// GetUserByID là phương thức nghiệp vụ để lấy thông tin người dùng bằng ID.
func (s *AccountService) GetUserByID(ctx context.Context, userID string) (*auth.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		// Kiểm tra lỗi "not found" từ repository (dựa vào thông báo lỗi từ repo)
		if strings.Contains(err.Error(), "not found") {
			return nil, fmt.Errorf("%w: ID %s", ErrUserNotFound, userID) // Dịch thành lỗi nghiệp vụ
		}
		return nil, fmt.Errorf("%w: failed to get user by ID: %v", ErrUserNotFound, err) // Wrap lỗi gốc
	}
	return user, nil
}

// UpdateUser là phương thức nghiệp vụ để cập nhật thông tin người dùng.
func (s *AccountService) UpdateUser(ctx context.Context, user *auth.User) error {
	// Nếu có mật khẩu mới được cung cấp trong struct user, băm nó trước khi cập nhật.
	// user.PasswordHash lúc này là mật khẩu thô được truyền từ handler/client
	if user.PasswordHash != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("%w: failed to hash password for update: %v", ErrHashPassword, err)
		}
		user.PasswordHash = string(hashedPassword) // Gán mật khẩu đã băm
	}

	err := s.userRepo.UpdateUser(ctx, user)
	if err != nil {
		// Xử lý các lỗi cụ thể từ repo nếu cần (ví dụ: duplicate email khi cập nhật)
		var dbxErr *dbx.DBXError
		if errors.As(err, &dbxErr) {
			if dbxErr.Code == dbx.DBXErrorCodeDuplicate {
				for _, field := range dbxErr.Fields {
					if strings.EqualFold(field, "Email") { // Sử dụng EqualFold
						return fmt.Errorf("%w: email already in use during update", ErrUpdateUserFail)
					}
				}
			}
			return fmt.Errorf("%w: dbx error: %s", ErrUpdateUserFail, dbxErr.Message)
		}
		return fmt.Errorf("%w: %v", ErrUpdateUserFail, err) // Wrap lỗi gốc
	}
	return nil
}

// DeleteUser là phương thức nghiệp vụ để xóa người dùng.
func (s *AccountService) DeleteUser(ctx context.Context, userID string) error {
	err := s.userRepo.DeleteUser(ctx, userID)
	if err != nil {
		// Kiểm tra lỗi "not found" từ repository (dựa vào thông báo lỗi từ repo)
		if strings.Contains(err.Error(), "not found") {
			return fmt.Errorf("%w: user with ID %s not found for deletion", ErrUserNotFound, userID)
		}
		// Xử lý lỗi ràng buộc tham chiếu nếu cần (ví dụ: người dùng có đơn hàng)
		var dbxErr *dbx.DBXError
		if errors.As(err, &dbxErr) {
			if dbxErr.Code == dbx.DBXErrorCodeReferenceConstraint {
				return fmt.Errorf("%w: cannot delete user due to data dependencies", ErrDeleteUserFail)
			}
		}
		return fmt.Errorf("%w: %v", ErrDeleteUserFail, err) // Wrap lỗi gốc
	}
	return nil
}
