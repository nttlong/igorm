package vexample

import (
	"context"
	"testing"
	"time"
	_ "vauth/models"
	"vcache"
	_ "vcache"
	"vdb"

	"github.com/stretchr/testify/assert"
)

var db *vdb.TenantDB

func initDb(driver string, conn string) error {
	_db, err := vdb.Open(driver, conn)
	if err != nil {
		return err
	}
	db = _db
	return nil
}
func TestInitMySqlDb(t *testing.T) {
	err := initDb("mysql", "root:123456@tcp(127.0.0.1:3306)/?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true")
	assert.NoError(t, err)
	tenantDb, err := db.CreateDB("a001")
	assert.NoError(t, err)
	assert.NotNil(t, tenantDb)
}
func TestNintMssql(t *testing.T) {
	err := initDb("sqlserver", "sqlserver://sa:123456@localhost:1433")
	assert.NoError(t, err)
	tenantDb, err := db.CreateDB("a001")
	assert.NoError(t, err)
	assert.NotNil(t, tenantDb)
}
func TestInitPostgres(t *testing.T) {
	err := initDb("postgres", "postgres://postgres:123456@localhost:5432?sslmode=disable")
	assert.NoError(t, err)
	tenantDb, err := db.CreateDB("a001")
	assert.NoError(t, err)
	assert.NotNil(t, tenantDb)
}
func TestVCacheInMemoryCache(t *testing.T) {
	/*
	 pc
	*/
	var cache vcache.Cache
	cache = vcache.NewInMemoryCache(10*time.Second, 10*time.Second)
	cache.Set(context.Background(), "key", "value", 10*time.Second)
}
func TestVCacheBaggerCache(t *testing.T) {

	var cache vcache.Cache
	cache, err := vcache.NewBadgerCache("./cache.db", "vcache")
	assert.NoError(t, err)
	cache.Set(context.Background(), "key", "value", 10*time.Second)
}
