package test

import (
	"testing"
	"vdb"
	"xauth/config"
	"xauth/repo"
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
	userRepo := repo.NewUserRepoSql(db)

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
		userRepo := repo.NewUserRepoSql(db)

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
	userRepo := repo.NewUserRepoSql(db)
	user, err := userRepo.GetUserById("00000000-0000-0000-0000-000000000000")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "admin", user.Username)
}
func TestLogin(t *testing.T) {
	cfg, err := config.NewConfig("./../config/config.yaml")
	assert.NoError(t, err)
	db, err := vdb.Open(cfg.Database.Driver, cfg.Database.Dsn)
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, err)
	assert.NotNil(t, db)
	var login services.Login
	login = services.NewLonginService(repo.NewUserRepoSql(db), services.NewAuthServiceArgon())
	user, err := login.DoLogin("admin", "123456")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "admin", user.Username)
}
func BenchmarkLogin(t *testing.B) {
	for i := 0; i < t.N; i++ {
		cfg, _ := config.NewConfig("./../config/config.yaml")

		db, _ := vdb.Open(cfg.Database.Driver, cfg.Database.Dsn)

		var login services.Login
		login = services.NewLonginService(repo.NewUserRepoSql(db), services.NewAuthServiceArgon())
		user, err := login.DoLogin("admin", "123456")
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "admin", user.Username)
	}
}
func BenchmarkLogin2(b *testing.B) {
	// Setup 1 lần duy nhất
	cfg, _ := config.NewConfig("./../config/config.yaml")
	db, _ := vdb.Open(cfg.Database.Driver, cfg.Database.Dsn)

	login := services.NewLonginService(
		repo.NewUserRepoSql(db),
		services.NewAuthServiceArgon(),
	)

	b.ResetTimer() // reset bộ đếm để bỏ qua thời gian setup

	for i := 0; i < b.N; i++ {
		user, err := login.DoLogin("admin", "123456")
		if err != nil {
			b.Fatalf("login failed: %v", err)
		}
		if user == nil || user.Username != "admin" {
			b.Fatalf("unexpected user: %+v", user)
		}
	}
}
func BenchmarkArgonHash(b *testing.B) {
	auth := services.NewAuthServiceArgon()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := auth.HashPassword("123456")
		if err != nil {
			b.Fatal(err)
		}
	}
}
func BenchmarkArgonVerify(b *testing.B) {
	auth := services.NewAuthServiceArgon()
	hash, _ := auth.HashPassword("123456")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ok, _ := auth.VerifyPassword(hash, "123456")
		if !ok {
			b.Fatal("verify failed")
		}
	}
}
