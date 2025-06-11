package auth

import (
	authModel "dbmodels/auth"
	"dbx"
	"strings"
)

func (u *User) Login(username string, password string) (*OAuth2Token, error) {
	CreateSysAdminUser(u.TenantDb, u.Context)
	// get user from db
	user, err := dbx.Query[authModel.User](u.TenantDb, u.Context).Where("Username = ?", username).First()

	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}

	// hashPassword, err := hashPasswordWithSalt(password + "@" + strings.ToLower(username))
	// if err != nil {
	// 	return false, err
	// }

	if err = verifyPassword(password+"@"+strings.ToLower(username), user.PasswordHash); err != nil {
		return nil, nil
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
