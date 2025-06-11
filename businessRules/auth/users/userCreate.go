package auth

import (
	authModels "dbmodels/auth"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (u *User) Create(username, password, email string) (*authModels.User, error) {
	hasPass, err := hashPasswordWithSalt(password + "@" + strings.ToLower(username))
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	// Tạo mới user
	user := &authModels.User{
		UserId:       uuid.New().String(),
		Username:     username,
		PasswordHash: hasPass,
		Email:        email,
		CreatedBy:    "system",
		CreatedAt:    time.Now().UTC(),
	}

	// Lưu vào cơ sở dữ liệu
	err = u.TenantDb.InsertWithContext(u.Context, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Trả về user mới
	return user, nil
}
