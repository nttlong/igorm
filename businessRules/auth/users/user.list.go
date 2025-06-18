package auth

import (
	authModels "dbmodels/auth"
	"dbx"
)

func (u *User) List(filter struct{}) (interface{}, error) {
	tokenInfo, err := u.ValidateAccessToken(u.AccessToken)
	if err != nil || tokenInfo == nil {

		return nil, err
	}
	ret, err := dbx.Query[authModels.User](u.TenantDb, u.Context).Items()
	return ret, err

}
