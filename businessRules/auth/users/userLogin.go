package auth

import (
	authModel "dbmodels/auth"
	"dbx"
	"dynacall"
	"fmt"

	service "unvs.br.auth/services"

	authError "unvs.br.auth/errors"
)

func (u *User) AuthenticateUser(username string, password string) (*service.OAuth2Token, error) {
	(&u.TokenService).DecodeAccessToken("token")
	err := CreateSysAdminUser(u.TenantDb, u.Context)
	if err != nil {
		return nil, err
	}
	if username == "" || password == "" {
		return nil, fmt.Errorf("username or password is empty")
	}
	// get user from db
	var user authModel.User

	if !u.Cache.Get(u.Context, "user_"+username, &user) {
		_user, err := dbx.Query[authModel.User](
			u.TenantDb,
			u.Context,
		).Select("Email,UserId,Id, Username, PasswordHash, RoleId").Where(
			"Username = ?",
			username).First()

		if err != nil {
			return nil, err
		} else {
			if _user == nil {
				return nil, &authError.AuthError{
					Code:    authError.ErrInvalidUsernameOrPassword,
					Message: "Invalid username or password",
				}
			}
			user = *_user
			u.Cache.Set(u.Context, "user_"+username, *_user, 0)

		}
	}

	key := fmt.Sprintf("user_%s %s", username, password)
	ok := ""
	if !u.Cache.Get(u.Context, key, &ok) {

		err := u.VerifyPassword(username, password, user.PasswordHash)
		if err != nil {
			return nil, &authError.AuthError{
				Code:    authError.ErrInvalidUsernameOrPassword,
				Message: "Invalid username or password",
			}
		}
		u.Cache.Set(u.Context, key, "ok", 0)
	}

	defaultRole := "user"
	if user.RoleId != nil {
		role, err := dbx.Query[authModel.Role](u.TenantDb, u.Context).Where("Id = ?", user.RoleId).First()
		if err != nil {
			return nil, err
		}
		if role != nil {
			defaultRole = role.Name
		}
	}

	return u.GenerateToken(struct {
		UserId   string
		RoleId   string
		Username string
		Email    *string
	}{
		UserId:   user.UserId,
		RoleId:   defaultRole,
		Username: user.Username,
		Email:    user.Email,
	})

}
func (u *User) Login(login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}) (*service.OAuth2Token, error) {
	ret, err := u.AuthenticateUser(login.Username, login.Password)
	if err != nil {
		if authErr, ok := err.(*authError.AuthError); ok {
			if authErr.Code == authError.ErrInvalidUsernameOrPassword {
				return nil, &dynacall.CallError{
					Code: dynacall.CallErrorCodeAuthenticationFailed,
					Err:  err,
				}
			}
			return nil, authErr
		}

		return nil, err
	}
	return ret, nil
}
