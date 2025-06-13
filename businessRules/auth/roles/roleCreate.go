package auth

import (
	_ "dbmodels/auth"
	dbmodels "dbmodels/auth"
	"dbx"
	"dynacall"
	"errors"
	"time"

	authErr "unvs.br.auth/errors"
)

func (r *RoleService) Create(data *struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}) (*dbmodels.Role, error) {
	tokenInfo, err := r.ValidateAccessToken(r.AccessToken)
	if err != nil {
		if auErr, ok := err.(*authErr.AuthError); ok {
			if auErr.Code == authErr.ErrTokenExpired {
				return nil, &dynacall.CallError{
					Code: dynacall.CallErrorCodeTokenExpired,
					Err:  err,
				}
			}
			return nil, &dynacall.CallError{
				Code: dynacall.CallErrorCodeAccessDenied,
				Err:  err,
			}
		}

		return nil, err
	}
	if tokenInfo == nil {
		return nil, &dynacall.CallError{
			Code: dynacall.CallErrorCodeAccessDenied,
			Err:  errors.New("Access Deny"),
		}
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
