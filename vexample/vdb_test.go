package vexample

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"
	_ "vauth/models"
	models "vauth/models"
	"vcache"
	_ "vcache"
	"vdb"

	"github.com/google/uuid"
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
func TestInitMssql(t *testing.T) {
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
func TestCacheRedis(t *testing.T) {
	var cache vcache.Cache
	cache = vcache.NewRedisCache("localhost:6379", "123456", "vcache", 0, 10*time.Second)

	cache.Set(context.Background(), "key", "value", 10*time.Second)
}
func TestCacheMemCached(t *testing.T) {
	var cache vcache.Cache
	cache = vcache.NewMemcachedCache("localhost:11211", "vcache")
	cache.Set(context.Background(), "key", "value", 10*time.Second)
}
func TestCreateUserMysql(t *testing.T) {
	vdb.SetManagerDb("mysql", "sys")
	err := initDb("mysql", "root:123456@tcp(127.0.0.1:3306)/sys?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true")
	assert.NoError(t, err)
	tenantDb, err := db.CreateDB("a002")
	assert.NoError(t, err)
	assert.NotNil(t, tenantDb)
	user := models.User{}
	err = tenantDb.First(&user, "id = ?", 1)
	if _, ok := err.(*vdb.ErrRecordNotFound); ok {
		user.Email = "test@test.com"
		user.UserId = uuid.NewString()
		user.Username = "test"
		user.HashPassword = "123456"
		user.IsActive = true
		user.CreatedAt = time.Now()
		err = tenantDb.Create(&user)
		assert.NoError(t, err)

	} else {
		assert.NoError(t, err)
	}
	assert.NoError(t, err)

}
func TestUpdateUserMysql(t *testing.T) {
	vdb.SetManagerDb("mysql", "sys")
	err := initDb("mysql", "root:123456@tcp(127.0.0.1:3306)/sys?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true")
	assert.NoError(t, err)
	tenantDb, err := db.CreateDB("a002")
	assert.NoError(t, err)
	assert.NotNil(t, tenantDb)
	user := models.User{}
	err = tenantDb.First(&user, "id = ?", 1)
	if _, ok := err.(*vdb.ErrRecordNotFound); ok {
		fmt.Print("not found")
	} else {
		assert.NoError(t, err)
	}
	for j := 0; j < 5; j++ {

		err = tenantDb.Model(&user).Where("id = ?", 2).Update("UserName", "NewName").Error
		if err != nil {
			log.Fatal(err)
		}
	}

}
func BenchmarkTestUpdateUserMysql(t *testing.B) {
	vdb.SetManagerDb("mysql", "sys")
	err := initDb("mysql", "root:123456@tcp(127.0.0.1:3306)/sys?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true")
	assert.NoError(t, err)
	tenantDb, err := db.CreateDB("a002")
	assert.NoError(t, err)
	assert.NotNil(t, tenantDb)

	for i := 0; i < t.N; i++ {

		err = tenantDb.Model(&models.User{}).Where("id = ?", 1).Update("UserName", "NewName11").Error
		assert.NoError(t, err)
	}

}
func BenchmarkTestCreateUserMysql(t *testing.B) {
	vdb.SetManagerDb("mysql", "sys")
	err := initDb("mysql", "root:123456@tcp(127.0.0.1:3306)/sys?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true")
	assert.NoError(t, err)
	tenantDb, err := db.CreateDB("a002")
	assert.NoError(t, err)
	assert.NotNil(t, tenantDb)
	for i := 0; i < t.N; i++ {
		for j := 0; j < 10000; j++ {
			user := models.User{}

			err = tenantDb.Where("email = ?", "test"+strconv.Itoa(i*100+j)+"@test.com").First(&user)
			if _, ok := err.(*vdb.ErrRecordNotFound); ok {
				user.Email = "test" + strconv.Itoa(i*100+j) + "@test.com"
				user.Username = "test" + strconv.Itoa(i*100+j)
				user.HashPassword = "123456"
				user.UserId = uuid.NewString()
				user.IsActive = true
				user.CreatedAt = time.Now()
				err = tenantDb.Create(&user)
				assert.NoError(t, err)

			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, err)
		}
	}
}
func TestDeleteUserMysql(t *testing.T) {
	vdb.SetManagerDb("mysql", "sys")
	err := initDb("mysql", "root:123456@tcp(127.0.0.1:3306)/sys?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true")
	assert.NoError(t, err)
	tenantDb, err := db.CreateDB("a002")
	assert.NoError(t, err)
	assert.NotNil(t, tenantDb)
	user := models.User{}
	err = tenantDb.First(&user, "id > ?", 1)
	if _, ok := err.(*vdb.ErrRecordNotFound); ok {
		fmt.Print("not found")
	} else {
		err = tenantDb.Delete(&user).Error
		assert.NoError(t, err)
	}

}
func BenchmarkTestDeleteUserMysql(t *testing.B) {
	vdb.SetManagerDb("mysql", "sys")
	err := initDb("mysql", "root:123456@tcp(127.0.0.1:3306)/sys?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true")
	assert.NoError(t, err)
	tenantDb, err := db.CreateDB("a002")
	assert.NoError(t, err)
	assert.NotNil(t, tenantDb)
	for i := 0; i < t.N; i++ {
		user := models.User{}

		err = tenantDb.Where("id = ?", 1).Delete(&user).Error
		assert.NoError(t, err)
	}

}
func BenchmarkTestGetThenDeleteUserMysql(b *testing.B) {
	vdb.SetManagerDb("mysql", "sys")
	err := initDb("mysql", "root:123456@tcp(127.0.0.1:3306)/sys?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true")
	assert.NoError(b, err)
	tenantDb, err := db.CreateDB("a002")
	assert.NoError(b, err)
	assert.NotNil(b, tenantDb)
	user := models.User{}
	for i := 0; i < b.N; i++ {
		err = tenantDb.First(&user, "id > ?", 1)
		if _, ok := err.(*vdb.ErrRecordNotFound); ok {
			fmt.Print("not found")
		} else {
			err = tenantDb.Delete(&user).Error
			assert.NoError(b, err)
		}
	}

}
