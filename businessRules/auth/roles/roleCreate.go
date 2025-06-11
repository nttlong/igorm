package auth

import (
	_ "dbmodels/auth"
	dbmodels "dbmodels/auth"
	"dbx"
	"time"
)

func (r *RoleService) Create(code, name, description string) (*dbmodels.Role, error) {
	if err := r.validateAccessToken(); err != nil {
		return nil, err
	}
	role := dbmodels.Role{
		RoleId:      dbx.NewUUID(),
		Code:        code,
		Name:        name,
		Description: dbx.FullTextSearchColumn(description),
		CreatedAt:   time.Now().UTC(),
		CreatedBy:   "system",
	}

	err := dbx.InsertWithContext(r.Context, r.TenantDb, &role)
	if err != nil {
		return nil, err
	}
	return &role, nil
}
