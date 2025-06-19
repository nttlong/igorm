/*
This file support function create sys admin user root
*/
package auth

import (
	"context"
	authModels "dbmodels/auth"
	"dbx"
	"sync"
	"time"

	"github.com/google/uuid"
	"unvs.br.auth/services"
)

var cacheCreatedSysAdminUser = sync.Map{}

func CreateSysAdminUser(db *dbx.DBXTenant, ctx context.Context) error {
	if _, ok := cacheCreatedSysAdminUser.Load(db.TenantDbName); ok {
		return nil
	}
	ret := createSysAdminUserNoCache(db, ctx)
	if ret != nil {
		return ret
	}
	cacheCreatedSysAdminUser.Store(db.TenantDbName, true)
	return nil
}

func createSysAdminUserNoCache(db *dbx.DBXTenant, ctx context.Context) error {

	c, err := dbx.CountWithContext[authModels.User](ctx, db, "Username = ?", "root")
	if err != nil {
		return err
	}
	if c > 0 {
		return nil
	}
	rootUser := &authModels.User{
		UserId:       uuid.New().String(),
		Username:     "root",
		PasswordHash: "root",
		CreatedBy:    "root",
		CreatedAt:    time.Now().UTC(),
		IsSupperUser: true,
		IsLocked:     false,
	}
	rootUser.PasswordHash, err = (&services.PasswordService{}).HashPassword("root", "root")
	if err != nil {
		return err
	}
	err = db.InsertWithContext(ctx, rootUser)
	if err != nil {
		return err
	}
	return nil
}
