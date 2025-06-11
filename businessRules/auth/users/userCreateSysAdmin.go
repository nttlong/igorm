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
)

var cacheCreatedSysAdminUser = sync.Map{}

func CreateSysAdminUser(db *dbx.DBXTenant, ctx context.Context) {
	if _, ok := cacheCreatedSysAdminUser.Load(db.TenantDbName); ok {
		return
	}
	createSysAdminUserNoCache(db, ctx)
	cacheCreatedSysAdminUser.Store(db.TenantDbName, true)
}

func createSysAdminUserNoCache(db *dbx.DBXTenant, ctx context.Context) {

	c, err := dbx.CountWithContext[authModels.User](ctx, db, "username = ?", "root")
	if err != nil {
		return
	}
	if c > 0 {
		return
	}
	rootUser := &authModels.User{
		UserId:       uuid.New().String(),
		Username:     "root",
		PasswordHash: "root",
		Email:        "root@test.com",
		CreatedBy:    "system",
		CreatedAt:    time.Now().UTC(),
		IsSupperUser: true,
		IsLocked:     false,
	}
	rootUser.PasswordHash, err = hashPasswordWithSalt("root@root")
	if err != nil {
		return
	}
	db.InsertWithContext(ctx, rootUser)
}
