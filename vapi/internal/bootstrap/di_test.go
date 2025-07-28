package bootstrap

import (
	"context"
	"testing"
	"time"
	accModels "vapi/internal/account/models"
	"vapi/internal/security/models"
	"vdb"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func BenchmarkTestDI(b *testing.B) {

	c := GetAppContainer("../config/config.yaml")
	assert.NoError(b, c.Error)
	assert.IsType(b, &AppContainer{}, c)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.GetContext = func() context.Context {
			return b.Context()
		}
		c.GetTenantName = func() string {
			return "tenant1"
		}

		c.Security.Get().CreateOrUpdate(&models.SecurityPolicy{
			TenantID:         "tenant1",
			JwtSecret:        c.SharedSecret.Get().Generate(),
			MaxLoginFailures: 5,
			LockoutMinutes:   15,
			JwtExpireMinutes: 60,
			CreatedAt:        time.Now().UTC(),
		})
	}

}
func BenchmarkTestDIGetPolicy(b *testing.B) {
	c := GetAppContainer("../config/config.yaml")
	assert.NoError(b, c.Error)
	assert.IsType(b, &AppContainer{}, c)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.GetContext = func() context.Context {
			return b.Context()
		}
		c.GetTenantName = func() string {
			return "tenant1"
		}

		_, err := c.Security.Get().Get()
		assert.NoError(b, err)

	}

}
func BenchmarkCreateAccount(b *testing.B) {
	c := GetAppContainer("../../cmd/config.yaml")
	assert.NoError(b, c.Error)
	assert.IsType(b, &AppContainer{}, c)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.GetContext = func() context.Context {
			return b.Context()
		}
		c.GetTenantName = func() string {
			return "tenant1"
		}
		pw, err := c.PwdSvc.Get().Hash("testpassword")
		assert.NoError(b, err)
		acc := &accModels.Account{
			UserID:         uuid.NewString(),
			Username:       "testuser",
			HashedPassword: pw,
			FullName:       "Test User",
			Email:          vdb.Ptr("test@test.com"),
			Role:           "user",
			CreatedAt:      time.Now().UTC(),
		}
		c.AccountSvc.Get().CreateOrUpdate(acc)
		assert.NoError(b, c.Error)
		// _, err := c.Security.Get().Get()
		// assert.NoError(b, err)

	}

}
func BenchmarkTestCheckPassword(b *testing.B) {
	c := GetAppContainer("../../cmd/config.yaml")
	assert.NoError(b, c.Error)
	assert.IsType(b, &AppContainer{}, c)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.GetContext = func() context.Context {
			return b.Context()
		}
		c.GetTenantName = func() string {
			return "tenant1"
		}

		pw, err := c.PwdSvc.Get().Hash("testpassword")

		assert.NoError(b, err)
		acc := &accModels.Account{
			UserID:         uuid.NewString(),
			Username:       "testuser",
			HashedPassword: pw,
			FullName:       "Test User",
			Email:          vdb.Ptr("test@test.com"),
			Role:           "user",
			CreatedAt:      time.Now().UTC(),
		}
		assert.NotEmpty(b, acc.HashedPassword)

	}
}
func BenchmarkTestLogin(b *testing.B) {
	c := GetAppContainer("../../cmd/config.yaml")
	assert.NoError(b, c.Error)
	assert.IsType(b, &AppContainer{}, c)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.GetContext = func() context.Context {
			return b.Context()
		}
		c.GetTenantName = func() string {
			tenantName := "tenant1"

			return tenantName
		}
		// c.GetDb = func() *vdb.TenantDB {
		// 	ret, err := c.Db.Get().CreateDB("tenant1")

		// 	assert.NoError(b, err)
		// 	return ret
		// }

		r, err := c.AccountSvc.Get().Login("testuser", "testpassword")
		assert.NoError(b, err)
		assert.NotNil(b, r)

	}
}
