package role

import (
	"dbx"
	"testing"
	"time"
	"unvs/internal/model/auth"

	"github.com/stretchr/testify/assert"
)

var TenantDb *dbx.DBXTenant

func TestCreateMssqlTenantDb(t *testing.T) {
	// arrange
	mssqlCfg := dbx.Cfg{
		Driver: "mssql",
		Host:   "localhost",

		User:     "sa",
		Password: "123456",
		SSL:      false,
	}
	db := dbx.NewDBX(mssqlCfg)
	db.Open()
	defer db.Close()
	err := db.Ping()
	assert.NoError(t, err)
	tenantDb, err := db.GetTenant("TenantDb")

	assert.NoError(t, err)
	TenantDb = tenantDb
	assert.NoError(t, err)
	TenantDb.Open()
	defer TenantDb.Close()
	err = TenantDb.Ping()
	assert.NoError(t, err)
}
func TestRoleRepo_Create(t *testing.T) {
	TestCreateMssqlTenantDb(t)
	repo := NewRoleRepository(TenantDb)
	// act
	role := &auth.Role{
		Code:        "admin",
		Name:        "Admin",
		CreatedBy:   "admin",
		Description: "This is the admin role",
		CreatedAt:   time.Now(),
	}
	TenantDb.Open()
	defer TenantDb.Close()
	err := repo.Create(nil, role)
	// assert
	assert.NoError(t, err)
}
