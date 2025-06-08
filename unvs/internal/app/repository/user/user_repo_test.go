package user

import (
	"context"
	"dbx"
	"fmt"
	_ "fmt"
	"testing"
	"time"

	"unvs/internal/model/auth"
	_ "unvs/internal/model/auth"
	_ "unvs/internal/model/base"

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

var DbTenant *dbx.DBXTenant

func TestCreateDbxTenant(t *testing.T) {
	db := dbx.NewDBX(createMssqlConfig())

	dbTenant, err := db.GetTenant("tenant1")
	assert.NoError(t, err, "Error creating tenant")

	DbTenant = dbTenant
}
func TestUserRepo_Create(t *testing.T) {
	TestCreateDbxTenant(t)
	DbTenant.Open()
	defer DbTenant.Close()
	repo := NewUserRepo(DbTenant)
	for i := 0; i < 1000; i++ {
		user := auth.User{
			Username:     "testuser",
			PasswordHash: "testpassword",
			Email:        "testemail",
			CreatedBy:    "testuser",

			Description: "Chỉ là test thôi",
			CreatedAt:   time.Now(),
		}

		start := time.Now()

		err := repo.CreateUser(context.Background(), &user)
		n := time.Since(start).Milliseconds()
		fmt.Println(n)
		if dbxErr, ok := err.(*dbx.DBXError); ok {
			if dbxErr.Code == dbx.DBXErrorCodeDuplicate {
				if dbxErr.Fields[0] == "email" {

					fmt.Println("Duplicate email")
					continue
				}
				if dbxErr.Fields[0] == "username" {
					fmt.Println("Duplicate username")
					continue
				}
			}
			t.Log(dbxErr.Message)

		} else {
			assert.NoError(t, err)
		}
	}

}
