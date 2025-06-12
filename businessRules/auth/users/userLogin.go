package auth

import (
	authModel "dbmodels/auth"
	"dbx"
	"fmt"
	"strings"
	"time"

	authError "unvs.br.auth/errors"
)

func (u *User) Login(username string, password string, LoginCount time.Time) (*OAuth2Token, error) {
	CreateSysAdminUser(u.TenantDb, u.Context)
	if username == "" || password == "" {
		return nil, fmt.Errorf("username or password is empty")
	}
	// get user from db
	var user authModel.User
	var err error
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
			user = *_user
			u.Cache.Set(u.Context, "user_"+username, *_user, 0)

		}
	}

	key := fmt.Sprintf("user_%s %s", username, user.PasswordHash)
	ok := ""
	if !u.Cache.Get(u.Context, key, &ok) {

		if err = verifyPassword(password+"@"+strings.ToLower(username), user.PasswordHash); err != nil {
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

	return generateToken(u.JwtSecret, user.UserId, defaultRole)

}
func (u *User) Login2(login struct {
	Username string    `json:"username"`
	Password string    `json:"password"`
	LoginOn  time.Time `json:"loginOn"`
}) (*OAuth2Token, error) {
	return u.Login(login.Username, login.Password, login.LoginOn)
}
