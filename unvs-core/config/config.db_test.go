package config

import (
	"dbx"
	"fmt"
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDBConfig(t *testing.T) {
	LoadConfig("../config")
	cfg := CreateCfg()
	assert.NotNil(t, cfg)
	assert.Equal(t, AppConfigInstance.Database.Driver, cfg.Driver)
	assert.Equal(t, AppConfigInstance.Database.Host, cfg.Host)
	assert.Equal(t, AppConfigInstance.Database.Port, cfg.Port)
	assert.Equal(t, AppConfigInstance.Database.User, cfg.User)
	assert.Equal(t, AppConfigInstance.Database.Password, cfg.Password)
	assert.Equal(t, AppConfigInstance.Database.SSL, cfg.SSL)
	assert.Equal(t, AppConfigInstance.Database.Name, cfg.DbName)
	assert.Equal(t, AppConfigInstance.Database.IsMultiTenancy, cfg.IsMultiTenancy)
}
func TestConnectDB(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {

			fmt.Println(string(debug.Stack()))

		}
	}() // Gọi ngay lập tức hàm ẩn danh deferred
	defer dbx.CloseAll()
	LoadConfig("../config")
	cfg := CreateCfg()
	db := dbx.NewDBX(*cfg)
	db.Open()

	err := db.Ping()
	assert.NoError(t, err)
	tenantDb, err := db.GetTenant("tenant2")
	assert.NoError(t, err)
	assert.NotNil(t, tenantDb)
	err = tenantDb.Ping()
	assert.NoError(t, err)

}
