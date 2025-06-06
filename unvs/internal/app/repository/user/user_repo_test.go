package user

import (
	"context"
	"dbx"
	"testing"
	"time"

	"unvs/internal/model/auth"
	_ "unvs/internal/model/auth"
	"unvs/internal/model/base"
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

var DbTenant dbx.DBXTenant

func TestCreateDbxTenant(t *testing.T) {
	db := dbx.NewDBX(createMssqlConfig())

	dbTenant, err := db.GetTenant("tenant1")
	assert.NoError(t, err, "Error creating tenant")

	DbTenant = *dbTenant
}
func TestUserRepo_Create(t *testing.T) {
	TestCreateDbxTenant(t)
	DbTenant.Open()
	defer DbTenant.Close()
	repo := NewUserRepo(DbTenant)
	user := auth.User{
		Username:     "testuser",
		PasswordHash: "testpassword",
		Email:        "testemail",
		BaseModel: base.BaseModel{
			CreatedBy: "testuser",

			Description: "Chỉ là test thôi",
			CreatedAt:   time.Now(),
		},
	}
	err := repo.CreateUser(context.Background(), &user)
	assert.NoError(t, err)
}
