package account

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"dbx"
	"unvs/internal/app/cache"
	userRepo "unvs/internal/app/repository/user"

	"github.com/stretchr/testify/assert"
)

func createMssqlConfig() dbx.Cfg {
	return dbx.Cfg{
		Driver:   "mssql",
		Host:     "localhost",
		Port:     1433,
		User:     "sa",
		Password: "123456",
	}

}

var DbTenant dbx.DBXTenant

func TestCreateDbxTenant(t *testing.T) {
	db := dbx.NewDBX(createMssqlConfig())

	dbTenant, err := db.GetTenant("tenant1")
	assert.NoError(t, err, "Error creating tenant")

	DbTenant = *dbTenant
}
func CreateCache(ownerType reflect.Type) cache.Cache {
	return cache.NewInMemoryCache(
		ownerType,
		5*time.Minute,
		10*time.Minute,
	)
}
func TestCreateAccount(t *testing.T) {
	TestCreateDbxTenant(t)
	DbTenant.Open()
	defer DbTenant.Close()
	repo := userRepo.NewUserRepo(DbTenant)
	// Test case 1: Create account successfully
	account := NewAccountService(repo, CreateCache(reflect.TypeOf(repo)))
	for i := 0; i < 1000; i++ {
		start := time.Now()
		user, err := account.CreateAccount(context.Background(), "user3", "user1@example.com", "password1")
		n := time.Since(start).Milliseconds()
		fmt.Println("Time: ", n)
		if err != nil {
			t.Error(err)
		} else {
			assert.NoError(t, err, "Error creating account")
			assert.Equal(t, "user1", user.Username)
		}
	}
	// assert.Equal(t, "password1", account.Password)
	// assert.Equal(t, 0, account.FailedAttempts)
	// assert.Equal(t, false, account.IsLocked)
}
