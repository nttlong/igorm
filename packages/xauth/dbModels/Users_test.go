package dbmodels

import (
	"fmt"
	"testing"
	"vdb"
	"xconfig"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	cfg, err := xconfig.NewConfig("./../../../config.yaml")
	assert.NoError(t, err)

	vdb.SetManagerDb("postgres", cfg.Database.Postgres.Manager)
	db, err := vdb.Open("postgres", cfg.Database.Postgres.Dsn)
	assert.NoError(t, err)
	err = db.Ping()
	assert.NoError(t, err)
	tenantDb, err := db.CreateDB("test001")
	assert.NoError(t, err)
	fmt.Println(tenantDb)
	defer tenantDb.Close()
	err = tenantDb.Create(&Users{
		Username:       "admin",
		HashedPassword: "admin",
	})
	assert.NoError(t, err)

}
