package users

import (
	"time"
	userModel "unvs/internal/model/auth"
	"unvs/views"

	"dbx"

	"github.com/google/uuid"
)

type User struct {
	views.BaseView
}
type userInfo struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	Description string `json:"description"`
}
type Response struct {
	Data  interface{} `json:"data"`
	Error interface{} `json:"error"`
}

// this is business logic for creating user
func (v *User) Create(userInfo userInfo) (*Response, error) {
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
		CreatedBy:    v.Claim.Username,
		CreatedAt:    time.Now(),
		IsSupperUser: false,
		IsLocked:     false,
		Description:  dbx.FullTextSearchColumn(userInfo.Description),
	}
	// qr := dbx.Query[userModel.User](&v.DbTenant, v.Context)
	err = dbx.InsertWithContext(v.Context, &v.DbTenant, &user)
	if err != nil {
		return &Response{
			Data:  nil,
			Error: err,
		}, nil

	}

	return &Response{
		Data:  user,
		Error: nil,
	}, nil
}
func init() {
	views.AddView(&User{
		BaseView: views.BaseView{
			ViewPath: "auth/users",
			IsAuth:   true,
		},
	})
}
