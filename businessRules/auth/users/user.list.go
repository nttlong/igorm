package auth

import (
	authModels "dbmodels/auth"
	"dbx"
)

func (u *User) List(filter struct {
	PageIndex int
	PageSize  int
	Sort      string
}) (interface{}, error) {
	tokenInfo, err := u.ValidateAccessToken(u.AccessToken)
	if err != nil || tokenInfo == nil {

		return nil, err
	}
	qr := dbx.Pager[authModels.User](u.TenantDb, u.Context)

	qr.Select("UserId,Username,Email,CreatedAt,UpdatedAt,CreatedBy,UpdatedBy")
	if filter.PageSize == 0 {
		filter.PageSize = 50
	}
	if filter.Sort == "" {
		filter.Sort = "Id asc"
	}
	qr.Page(filter.PageIndex).Size(filter.PageSize).Sort(filter.Sort)

	ret, err := qr.Query()
	return ret, err

}
