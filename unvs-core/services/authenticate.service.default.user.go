package services

import (
	"context"
	"dbx"
	"sync"
	"time"

	"github.com/google/uuid"
	"unvs.core/models"
)

var cacheCreatedSysAdminUser = sync.Map{}

func (u *AuthenticateService) CreateSysAdminUser(db *dbx.DBXTenant, ctx context.Context) error {
	if _, ok := cacheCreatedSysAdminUser.Load(db.TenantDbName); ok {
		return nil
	}
	ret := u.createSysAdminUserNoCache(db, ctx)
	if ret != nil {
		return ret
	}
	cacheCreatedSysAdminUser.Store(db.TenantDbName, true)
	return nil
}
func (u *AuthenticateService) createSysAdminUserNoCache(db *dbx.DBXTenant, ctx context.Context) error {

	c, err := dbx.CountWithContext[models.User](ctx, db, "Username = ?", "root")
	if err != nil {
		return err
	}
	if c > 0 {
		return nil
	}
	rootUser := &models.User{
		UserId:       uuid.New().String(),
		Username:     "root",
		PasswordHash: "root",
		CreatedBy:    "root",
		CreatedAt:    time.Now().UTC(),
		IsSupperUser: true,
		IsLocked:     false,
	}
	rootUser.PasswordHash, err = (&PasswordService{}).HashPassword("root", "root")
	if err != nil {
		return err
	}
	err = db.InsertWithContext(ctx, rootUser)
	if err != nil {
		return err
	}
	return nil
}
