package role // Tên package là 'role', theo tên thư mục hoặc dịch vụ

import (
	"context"
	"fmt" // Để dùng trong các lỗi trả về hoặc debug

	"dbx"                      // Giả sử đây là package DBXTenant của bạn
	"unvs/internal/model/auth" // Import struct Role từ model/auth
)

// RoleRepository là interface định nghĩa các phương thức truy cập dữ liệu cho Role.
// Sử dụng interface giúp tách rời AccountService khỏi chi tiết triển khai DB.
type RoleRepository interface {
	Create(ctx context.Context, role *auth.Role) error
	GetByRoleId(ctx context.Context, roleId string) (*auth.Role, error)
	GetByCode(ctx context.Context, code string) (*auth.Role, error)
	Update(ctx context.Context, role *auth.Role) error
	DeleteByCode(ctx context.Context, code string) error
	DeleteByRoleId(ctx context.Context, roleId string) error
}

// dbxRoleRepository là một triển khai của RoleRepository sử dụng dbx.DBXTenant.
type dbxRoleRepository struct {
	dbxClient *dbx.DBXTenant // Client kết nối cơ sở dữ liệu
}

// NewRoleRepository tạo một instance mới của dbxRoleRepository.
// Nó nhận một client DBXTenant và trả về một triển khai của RoleRepository interface.
func NewRoleRepository(dbxClient *dbx.DBXTenant) RoleRepository {
	return &dbxRoleRepository{dbxClient: dbxClient}
}

// Create thêm một Role mới vào cơ sở dữ liệu.
func (r *dbxRoleRepository) Create(ctx context.Context, role *auth.Role) error {
	// Giả sử dbx.InsertWithContext có thể nhận bất kỳ struct nào
	return dbx.InsertWithContext(ctx, r.dbxClient, role)
}

// GetByRoleId lấy thông tin Role từ cơ sở dữ liệu bằng RoleId.
func (r *dbxRoleRepository) GetByRoleId(ctx context.Context, roleId string) (*auth.Role, error) {
	role, err := dbx.Query[auth.Role](r.dbxClient, ctx).Where("RoleId = ?", roleId).First()
	if err != nil {
		// Xử lý trường hợp không tìm thấy hoặc lỗi khác từ DBX
		if err == dbx.ErrNoRows { // Giả sử dbx có một lỗi cụ thể khi không tìm thấy hàng
			return nil, nil // Trả về nil Role, nil error nếu không tìm thấy
		}
		return nil, fmt.Errorf("failed to get role by ID: %w", err)
	}
	return role, nil
}

// GetByCode lấy thông tin Role từ cơ sở dữ liệu bằng Code.
func (r *dbxRoleRepository) GetByCode(ctx context.Context, code string) (*auth.Role, error) {
	role, err := dbx.Query[auth.Role](r.dbxClient, ctx).Where("Code = ?", code).First()
	if err != nil {
		// Xử lý trường hợp không tìm thấy hoặc lỗi khác từ DBX
		if err == dbx.ErrNoRows { // Giả sử dbx có một lỗi cụ thể khi không tìm thấy hàng
			return nil, nil // Trả về nil Role, nil error nếu không tìm thấy
		}
		return nil, fmt.Errorf("failed to get role by code: %w", err)
	}
	return role, nil
}

// Update cập nhật thông tin Role trong cơ sở dữ liệu.
func (r *dbxRoleRepository) Update(ctx context.Context, role *auth.Role) error {
	// Giả sử dbx.UpdateWithContext có thể nhận bất kỳ struct nào
	// và tự động xác định bản ghi để cập nhật dựa trên primary key (ví dụ: RoleId)
	return dbx.Query[auth.Role](r.dbxClient, ctx).Where("RoleId = ?", role.RoleId).Update(*role)

}

// DeleteByCode xóa Role khỏi cơ sở dữ liệu bằng Code.
func (r *dbxRoleRepository) DeleteByCode(ctx context.Context, code string) error {
	// Giả sử dbx.Query[auth.Role](...).Where(...).Delete() hoạt động như ví dụ User
	err := dbx.Query[auth.Role](r.dbxClient, ctx).Where("Code = ?", code).Delete()
	if err != nil {
		return fmt.Errorf("failed to delete role by code: %w", err)
	}
	return nil
}

// DeleteByRoleId xóa Role khỏi cơ sở dữ liệu bằng RoleId.
func (r *dbxRoleRepository) DeleteByRoleId(ctx context.Context, roleId string) error {
	// Giả sử dbx.Query[auth.Role](...).Where(...).Delete() hoạt động như ví dụ User
	err := dbx.Query[auth.Role](r.dbxClient, ctx).Where("RoleId = ?", roleId).Delete()
	if err != nil {
		return fmt.Errorf("failed to delete role by ID: %w", err)
	}
	return nil
}
