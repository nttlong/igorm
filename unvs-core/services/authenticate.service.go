package services

import (
	"context"
	"dbx"
	"fmt"
	"time"

	"github.com/google/uuid"
	cacher "unvs.core/cacher"
	coreErrors "unvs.core/errors"
	"unvs.core/models"
	_ "unvs.core/models"
)

type AuthenticateService struct {
	TokenService
	PasswordService
	Context context.Context
	Cache   cacher.Cache
}

func (u *AuthenticateService) AuthenticateUser(username string, password string) (*OAuth2Token, error) {
	(&u.TokenService).DecodeAccessToken("token")
	err := u.CreateSysAdminUser(u.TenantDb, u.Context)
	if err != nil {
		return nil, err
	}
	if username == "" || password == "" {
		return nil, fmt.Errorf("username or password is empty")
	}
	// get user from db
	var user models.User
	cacheKey := fmt.Sprintf("%s_user_%s", u.TenantDb.TenantDbName, username)
	if !u.Cache.Get(u.Context, cacheKey, &user) {
		_user, err := dbx.Query[models.User](
			u.TenantDb,
			u.Context,
		).Select("Email,UserId,Id, Username, PasswordHash, RoleId").Where(
			"Username = ?",
			username).First()

		if err != nil {
			return nil, err
		} else {
			if _user == nil {
				return nil, &coreErrors.CoreError{
					Code:    coreErrors.Error_LoginFailed,
					Message: "Invalid username or password",
				}
			}
			user = *_user
			u.Cache.Set(u.Context, "user_"+username, *_user, 0)

		}
	}

	key := fmt.Sprintf("%s_user_%s %s", u.TenantDb.TenantDbName, username, password)
	ok := ""
	if !u.Cache.Get(u.Context, key, &ok) {

		err := u.VerifyPassword(username, password, user.PasswordHash)
		if err != nil {
			return nil, &coreErrors.CoreError{
				Code:    coreErrors.Error_LoginFailed,
				Message: "Invalid username or password",
			}
		}

		u.Cache.Set(u.Context, key, "ok", 0)
	}

	defaultRole := "user"
	if user.RoleId != nil {
		role, err := dbx.Query[models.Role](u.TenantDb, u.Context).Where("Id = ?", user.RoleId).First()
		if err != nil {
			return nil, err
		}
		if role != nil {
			defaultRole = role.Name
		}
	}

	ret, err := u.GenerateToken(GenerateTokenParams{
		UserId:   user.UserId,
		RoleId:   defaultRole,
		Username: user.Username,
		Email:    user.Email,
	})
	if err != nil {
		return nil, err
	}
	info := &models.LoginInfo{
		Id:           uuid.New().String(),
		UserId:       ret.UserId,
		RefreshToken: ret.RefreshToken,
		AccessToken:  &ret.AccessToken,
		CreatedAt:    time.Now(),
	}
	u.producer(info)
	return ret, nil

}
