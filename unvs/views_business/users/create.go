package users

import (
	userModel "unvs/internal/model/auth"
	"unvs/views"

	"github.com/google/uuid"
)

type User struct {
	views.BaseView
}
type userInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// this is business logic for creating user
func (v *User) Create(userInfo userInfo) (*userModel.User, error) {
	// Code to create user
	hashPassword, err := v.hashPassword(userInfo)
	if err != nil {
		return nil, err
	}
	user := userModel.User{
		UserId:       uuid.New().String(),
		Username:     userInfo.Username,
		PasswordHash: hashPassword,
		Email:        userInfo.Email,
	}

	return &user, nil
}
func init() {
	views.AddView(&User{
		BaseView: views.BaseView{
			ViewPath: "auth/users",
			IsAuth:   true,
		},
	})
}
