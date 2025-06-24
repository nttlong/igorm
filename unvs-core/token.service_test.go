package unvscore

import (
	"dbx"
	"fmt"
	"testing"
	"time"

	config "unvs.core/config"
	"unvs.core/services"

	"github.com/stretchr/testify/assert"
)

var TenantDb *dbx.DBXTenant

func TestGenerateToken(t *testing.T) {

	config.LoadConfig("../config")
	cfg := config.CreateCfg()
	AppConfigInstance := config.AppConfigInstance
	assert.NotNil(t, cfg)
	assert.Equal(t, AppConfigInstance.Database.Driver, cfg.Driver)
	assert.Equal(t, AppConfigInstance.Database.Host, cfg.Host)
	assert.Equal(t, AppConfigInstance.Database.Port, cfg.Port)
	assert.Equal(t, AppConfigInstance.Database.User, cfg.User)
	assert.Equal(t, AppConfigInstance.Database.Password, cfg.Password)
	assert.Equal(t, AppConfigInstance.Database.SSL, cfg.SSL)
	assert.Equal(t, AppConfigInstance.Database.Name, cfg.DbName)
	assert.Equal(t, AppConfigInstance.Database.IsMultiTenancy, cfg.IsMultiTenancy)
	dbx := config.CreateDbx()
	err := dbx.Open()
	assert.NoError(t, err)
	defer dbx.Close()
	tenantDb, err := dbx.GetTenant("tenant2")
	assert.NoError(t, err)
	Config.LoadConfig("../config")
	tokenService, err := Factory.GetTokenService(
		t.Context(), "vi", "tenant2",
	)
	assert.NoError(t, err)
	ret, err := tokenService.GenerateToken(services.GenerateTokenParams{
		UserId:   "123456",
		RoleId:   "123456",
		Username: "test",
		Email:    nil,
	})
	assert.NoError(t, err)
	assert.NotNil(t, ret)
	TenantDb = tenantDb
}
func TestAuthenticate(t *testing.T) {
	TestGenerateToken(t)
	svc, err := Factory.GetAuthService(
		t.Context(),
		"vi",
		"tenant2",
	)
	assert.NoError(t, err)
	avg := int64(0)
	for i := 0; i < 10000; i++ {
		start := time.Now()
		r, e := svc.AuthenticateUser("root", "root")
		n := time.Since(start).Milliseconds()
		fmt.Print("AuthenticateUser took ", n, " ms\n")
		avg += n
		assert.NoError(t, e)
		assert.NotEmpty(t, r)
	}
	avg /= 10000
	fmt.Print("AuthenticateUser took ", avg, " ms\n")
}
