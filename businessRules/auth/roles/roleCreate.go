package auth

import (
	_ "dbmodels/auth"
	dbmodels "dbmodels/auth"
	"dbx"
	"time"

	authErrors "unvs.br.auth/errors"
)

type RoleInfo struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (r *RoleService) Create(data []*RoleInfo) (*dbmodels.Role, error) {
	tokenInfo, err := r.ValidateAccessToken(r.AccessToken)
	if err != nil {
		return nil, err
	}
	if tokenInfo == nil {
		return nil, &authErrors.AuthError{
			Code:    authErrors.ErrAccessDeny,
			Message: "Access Deny",
		}
	}
	role := dbmodels.Role{
		RoleId:      dbx.NewUUID(),
		Code:        data[0].Code,
		Name:        data[0].Name,
		Description: dbx.FullTextSearchColumn(data[0].Description),
		CreatedAt:   time.Now().UTC(),
		CreatedBy:   "system",
	}

	err = dbx.InsertWithContext(r.Context, r.TenantDb, &role)
	if err != nil {
		return nil, err
	}
	return &role, nil
}
