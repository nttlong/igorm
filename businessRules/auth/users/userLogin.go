package auth

import (
	authModel "dbmodels/auth"
	"dbx"
	"fmt"

	service "unvs.br.auth/services"

	authError "unvs.br.auth/errors"
)

func (u *User) Login(username string, password string) (*service.OAuth2Token, error) {
	(&u.TokenService).DecodeAccessToken("token")
	CreateSysAdminUser(u.TenantDb, u.Context)
	if username == "" || password == "" {
		return nil, fmt.Errorf("username or password is empty")
	}
	// get user from db
	var user authModel.User

	if !u.Cache.Get(u.Context, "user_"+username, &user) {
		_user, err := dbx.Query[authModel.User](
			u.TenantDb,
			u.Context,
		).Where(
			"Username = ?",
			username).Select("Id, Username, PasswordHash, RoleId").First()

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

	key := fmt.Sprintf("user_%s %s", username, user.PasswordHash)
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

	return u.GenerateToken(user.UserId, defaultRole)

}
func (u *User) Login2(login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}) (*service.OAuth2Token, error) {
	return u.Login(login.Username, login.Password)
}
