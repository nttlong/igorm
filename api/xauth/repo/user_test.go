package repo

import (
	"testing"
	"vdb"
	"xauth/config"
	"xauth/services"

	"github.com/stretchr/testify/assert"
)

func TestCreateDefautUser(t *testing.T) {
	var authSvc services.PasswordService
	authSvc = services.NewAuthServiceArgon()
	cfg, err := config.NewConfig("./../config/config.yaml")
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	db, err := vdb.Open(cfg.Database.Driver, cfg.Database.Dsn)
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, err)
	assert.NotNil(t, db)
	userRepo := NewUserRepoSql(db)

	assert.NotNil(t, userRepo)
	hashPassword, err := authSvc.HashPassword("admin@123456")
	if err != nil {
		t.Fatal(err)
	}
	err = userRepo.CreateDefaultlUser(hashPassword)
	assert.NoError(t, err)

}
func BenchmarkCreateDefautUser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var authSvc services.PasswordService
		authSvc = services.NewAuthServiceArgon()

		cfg, err := config.NewConfig("./../config/config.yaml")
		assert.NoError(b, err)
		if err != nil {
			b.Fatal(err)
		}
		assert.NotNil(b, cfg)
		db, err := vdb.Open(cfg.Database.Driver, cfg.Database.Dsn)
		if err != nil {
			b.Fatal(err)
		}
		defer db.Close()
		assert.NoError(b, err)
		assert.NotNil(b, db)
		userRepo := NewUserRepoSql(db)

		assert.NotNil(b, userRepo)
		hashPassword, err := authSvc.HashPassword("admin@123456")
		if err != nil {
			b.Fatal(err)
		}
		err = userRepo.CreateDefaultlUser(hashPassword)
		assert.NoError(b, err)
	}
}
func TestGetUser(t *testing.T) {
	cfg, err := config.NewConfig("./../config/config.yaml")
	assert.NoError(t, err)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, cfg)
	db, err := vdb.Open(cfg.Database.Driver, cfg.Database.Dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	userRepo := NewUserRepoSql(db)
	user, err := userRepo.GetUserById("00000000-0000-0000-0000-000000000000")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "admin", user.Username)
}
