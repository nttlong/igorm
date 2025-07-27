package bootstrap

import (
	"context"
	"fmt"
	"testing"
	"time"
	accModels "vapi/internal/account/models"
	"vapi/internal/security/models"
	"vdb"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
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
		bhashed, err := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)
		fmt.Println(string(bhashed))
		password := "testpassword"
		hashed := "$10$LuQUbi0j1c2aHACpVP5jH.5cZQGOOU.1ZeOl2ebFlZkZGX9FEo4Q."
		fmt.Println(len(string(bhashed)))
		fmt.Println(len(hashed))

		err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
		if err != nil {
			fmt.Println("❌ Không khớp:", err)
		} else {
			fmt.Println("✅ Khớp!")
		}
		v := c.PwdSvc.Get().Verify("testpassword", "$10$LuQUbi0j1c2aHACpVP5jH.5cZQGOOU.1ZeOl2ebFlZkZGX9FEo4Q.")
		assert.True(b, v)
		pw, err := c.PwdSvc.Get().Hash("testpassword")
		fmt.Println(pw)
		fmt.Println(len(pw))
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
		acc2 := &accModels.Account{}
		db, err := c.TenantDb.Get().Tenant("tenant1")
		err = db.First(acc2, "username = ?", acc.Username)
		fmt.Println(acc2.ID)
		assert.NoError(b, c.Error)
		v = c.PwdSvc.Get().Verify("testpassword", "$2a$10$TcQvnH/7r1yPJ0eh/lYhT.SV7yazy8shwHUbt8m1hdi2YKn0eA42y")
		// _, err := c.Security.Get().Get()
		assert.NoError(b, err)

	}
}
