package bootstrap

import (
	"testing"
	"vdb"

	"github.com/stretchr/testify/assert"
)

func TestDI(t *testing.T) {
	// Create a new container
	c, err := NewAppContainer()
	assert.NoError(t, err)
	assert.IsType(t, &AppContainer{}, c)
	// Get the config service
}
func BenchmarkTestDI(b *testing.B) {
	c, err := NewAppContainer()

	assert.NoError(b, err)
	assert.IsType(b, &AppContainer{}, c)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		db, err := c.Tenant.Get().Tenant("test0001")
		assert.NoError(b, err)
		assert.IsType(b, &vdb.TenantDB{}, db)
	}

}
