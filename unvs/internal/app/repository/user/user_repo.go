package user // Tên package là 'user', theo tên thư mục

import (
	"context"
	"database/sql"
	"fmt"

	_ "fmt"
	_ "time"

	_ "dbx"                    // Giả sử đây là package DBXTenant của bạn
	"unvs/internal/model/auth" // Import struct User từ model/auth

	"dbx"

	"github.com/google/uuid"
)

// UserRepository là interface định nghĩa các phương thức truy cập dữ liệu cho User.
// Sử dụng interface giúp tách rời AccountService khỏi chi tiết triển khai DB.
type UserRepository interface {
	CreateUser(ctx context.Context, user *auth.User) error
	GetUserByID(ctx context.Context, userID string) (*auth.User, error)
	GetUserByEmail(ctx context.Context, email string) (*auth.User, error)
	UpdateUser(ctx context.Context, user *auth.User) error
	DeleteUser(ctx context.Context, userID string) error
	// Có thể thêm các phương thức khác tùy theo nghiệp vụ
	// GetAllUsers(ctx context.Context, limit, offset int) ([]*auth.User, error)
}
type dbxUserRepository struct {
	dbxClient dbx.DBXTenant // Client kết nối cơ sở dữ liệu
}

// NewUserRepository tạo một instance mới của UserRepository.
// Nó nhận một client DBXTenant và trả về một triển khai của UserRepository interface.
func NewUserRepo(dbxClient dbx.DBXTenant) UserRepository {
	return &dbxUserRepository{dbxClient: dbxClient}
}

// Add new user to database.
// CreateUser thêm một người dùng mới vào cơ sở dữ liệu.
func (r *dbxUserRepository) CreateUser(ctx context.Context, user *auth.User) error {
	user.UserId = uuid.New().String()
	return dbx.InsertWithContext(ctx, &r.dbxClient, user)
}

// GetUserByID lấy thông tin người dùng từ cơ sở dữ liệu bằng ID.
func (r *dbxUserRepository) GetUserByID(ctx context.Context, userID string) (*auth.User, error) {
	query := `SELECT id, email, password_hash, created_at, updated_at FROM users WHERE id = ?`
	user := &auth.User{}

	// Sử dụng QueryRowContext để lấy một hàng duy nhất.
	err := r.dbxClient.QueryRowContext(ctx, query, userID).Scan(
		&user.Id,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("repository: user with ID %s not found", userID)
		}
		return nil, fmt.Errorf("repository: failed to get user by ID %s: %w", userID, err)
	}
	return user, nil
}

// GetUserByEmail lấy thông tin người dùng từ cơ sở dữ liệu bằng email.
func (r *dbxUserRepository) GetUserByEmail(ctx context.Context, email string) (*auth.User, error) {
	query := `SELECT id, email, password_hash, created_at, updated_at FROM users WHERE email = ?`
	user := &auth.User{}

	err := r.dbxClient.QueryRowContext(ctx, query, email).Scan(
		&user.Id,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("repository: user with email %s not found", email)
		}
		return nil, fmt.Errorf("repository: failed to get user by email %s: %w", email, err)
	}
	return user, nil
}

// UpdateUser cập nhật thông tin người dùng trong cơ sở dữ liệu.
func (r *dbxUserRepository) UpdateUser(ctx context.Context, user *auth.User) error {
	panic("implement me")
}

// DeleteUser xóa người dùng khỏi cơ sở dữ liệu bằng ID.
func (r *dbxUserRepository) DeleteUser(ctx context.Context, userID string) error {
	query := `DELETE FROM users WHERE id = ?`

	result, err := r.dbxClient.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("repository: failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository: failed to get rows affected after delete: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("repository: no user found with ID %s to delete", userID)
	}
	return nil
}
