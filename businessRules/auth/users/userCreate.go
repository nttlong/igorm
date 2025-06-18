package auth

import (
	authModels "dbmodels/auth"
	"fmt"
	"time"

	"unvs.br.auth/services"

	"github.com/google/uuid"
)

func (u *User) Create(data struct {
	Username string  `json:"username" validate:"required,alphanum,min=3,max=255"`
	Password string  `json:"password" validate:"required,min=8,max=255"`
	Email    *string `json:"email" validate:"required,email"`
}) (*authModels.User, error) {
	tokenInfo, err := u.ValidateAccessToken(u.AccessToken)
	if err != nil || tokenInfo == nil {

		return nil, err
	}
	hasPass, err := (&services.PasswordService{}).HashPassword(data.Password, data.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	// Tạo mới user
	user := &authModels.User{
		UserId:       uuid.New().String(),
		Username:     data.Username,
		PasswordHash: hasPass,
		Email:        data.Email,
		CreatedBy:    tokenInfo.Username,
		CreatedAt:    time.Now().UTC(),
	}

	// Lưu vào cơ sở dữ liệu
	err = u.TenantDb.InsertWithContext(u.Context, user)
	if err != nil {

		return nil, err
	}

	// Trả về user mới
	return user, nil
}
