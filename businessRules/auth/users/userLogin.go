package auth

import (
	authModel "dbmodels/auth"
	"dbx"
	"dynacall"
	"fmt"
	"time"

	"github.com/google/uuid"
	service "unvs.br.auth/services"

	authError "unvs.br.auth/errors"
)

type LoginInfo struct {
	dbx.EntityModel
	Id           string  `db:"varchar(36);pk"`
	UserId       string  `db:"varchar(36);idx"` // user id
	RefreshToken string  `db:"varchar(255);idx"`
	AccessToken  *string `db:"varchar(500);idx"`

	CreatedAt time.Time
}

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
	cacheKey := fmt.Sprintf("%s_user_%s", u.TenantDb.TenantDbName, username)
	if !u.Cache.Get(u.Context, cacheKey, &user) {
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

	key := fmt.Sprintf("%s_user_%s %s", u.TenantDb.TenantDbName, username, password)
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

	ret, err := u.GenerateToken(struct {
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
	if err != nil {
		return nil, err
	}
	info := &LoginInfo{
		Id:           uuid.New().String(),
		UserId:       ret.UserId,
		RefreshToken: ret.RefreshToken,
		AccessToken:  &ret.AccessToken,
		CreatedAt:    time.Now(),
	}
	u.producer(info)
	return ret, nil

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

type LoginChanItem struct {
	info *LoginInfo
	db   *dbx.DBXTenant
}

var loginChange = make(chan LoginChanItem, 1000)

func (u *User) producer(info *LoginInfo) {
	item := LoginChanItem{
		info: info,
		db:   u.TenantDb.Clone("loginInfoDB"),
	}
	select {
	case loginChange <- item:
	default:
		fmt.Println("loginChange channel is full, dropping item")
	}
}
func consumer() {
	for {
		item := <-loginChange
		err := item.db.Insert(item.info)

		if err != nil {
			fmt.Println("error inserting login info", err)
		}
		_, err = item.db.Update(&authModel.User{}).Where(
			"UserId = ?", item.info.UserId).Set(
			"LastLoginAt", time.Now(),
		).Execute()
		if err != nil {
			fmt.Println("error updating user last login at", err)
		} else {

			fmt.Println("updated user last login")
		}

	}
}

func init() {
	dbx.AddEntities(&LoginInfo{})
	go consumer()

}
