package auth

import (
	_ "dbmodels/auth"
	dbmodels "dbmodels/auth"
	"dbx"
	"time"
)

func (r *RoleService) Create(data struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}) (*dbmodels.Role, error) {
	tokenInfo, err := r.ValidateAccessToken(r.AccessToken)
	if err != nil || tokenInfo == nil {

		return nil, err
	}

	role := dbmodels.Role{
		RoleId:      dbx.NewUUID(),
		Code:        data.Code,
		Name:        data.Name,
		Description: dbx.FullTextSearchColumn(data.Description),
		CreatedAt:   time.Now().UTC(),
		CreatedBy:   "system",
	}

	err = dbx.InsertWithContext(r.Context, r.TenantDb, &role)
	if err != nil {
		return nil, err
	}
	return &role, nil
}
