package bootstrap

import (
	"context"
	"testing"
	"time"
	"vapi/internal/security/models"

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
