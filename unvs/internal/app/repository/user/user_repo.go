package user // Tên package là 'user', theo tên thư mục

import (
	"context"

	_ "fmt"
	_ "time"

	_ "dbx"                    // Giả sử đây là package DBXTenant của bạn
	"unvs/internal/model/auth" // Import struct User từ model/auth

	"dbx"
)

// UserRepository là interface định nghĩa các phương thức truy cập dữ liệu cho User.
// Sử dụng interface giúp tách rời AccountService khỏi chi tiết triển khai DB.
type UserRepository interface {
	CreateUser(ctx context.Context, user *auth.User) error
	GetUserByID(ctx context.Context, userID string) (*auth.User, error)
	GetUserByEmail(ctx context.Context, email string) (*auth.User, error)
	GetUserByUsername(ctx context.Context, username string) (*auth.User, error)
	UpdateUser(ctx context.Context, user *auth.User) error
	DeleteUser(ctx context.Context, userID string) error
	// Có thể thêm các phương thức khác tùy theo nghiệp vụ
	// GetAllUsers(ctx context.Context, limit, offset int) ([]*auth.User, error)
}
type dbxUserRepository struct {
	dbxClient *dbx.DBXTenant // Client kết nối cơ sở dữ liệu
}

// NewUserRepository tạo một instance mới của UserRepository.
// Nó nhận một client DBXTenant và trả về một triển khai của UserRepository interface.
func NewUserRepo(dbxClient *dbx.DBXTenant) UserRepository {

	return &dbxUserRepository{dbxClient: dbxClient}
}

// Add new user to database.
// CreateUser thêm một người dùng mới vào cơ sở dữ liệu.
func (r *dbxUserRepository) CreateUser(ctx context.Context, user *auth.User) error {

	return dbx.InsertWithContext(ctx, r.dbxClient, user)
}

// GetUserByID lấy thông tin người dùng từ cơ sở dữ liệu bằng ID.
func (r *dbxUserRepository) GetUserByID(ctx context.Context, userID string) (*auth.User, error) {

	user, err := dbx.Query[auth.User](r.dbxClient, ctx).Where("UserId = ?", userID).First()
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByEmail lấy thông tin người dùng từ cơ sở dữ liệu bằng email.
func (r *dbxUserRepository) GetUserByEmail(ctx context.Context, email string) (*auth.User, error) {
	user, err := dbx.Query[auth.User](r.dbxClient, ctx).Where("Email = ?", email).First()
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByEmail lấy thông tin người dùng từ cơ sở dữ liệu bằng email.
func (r *dbxUserRepository) GetUserByUsername(ctx context.Context, username string) (*auth.User, error) {
	user, err := dbx.Query[auth.User](r.dbxClient, ctx).Where("Username = ?", username).First()
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser cập nhật thông tin người dùng trong cơ sở dữ liệu.
func (r *dbxUserRepository) UpdateUser(ctx context.Context, user *auth.User) error {
	panic("implement me")
}

// DeleteUser xóa người dùng khỏi cơ sở dữ liệu bằng ID.
func (r *dbxUserRepository) DeleteUser(ctx context.Context, userID string) error {
	return dbx.Query[auth.User](r.dbxClient, ctx).Where("UserId = ?", userID).Delete()
}
