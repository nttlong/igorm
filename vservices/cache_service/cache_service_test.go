package cache_service

// import (
// 	"testing"
// 	"vapi/internal/config"

// 	"github.com/stretchr/testify/assert"
// )

// func TestCacheService(t *testing.T) {
// 	cfg, err := config.NewConfigService("./../config/config.yaml")
// 	assert.NoError(t, err)
// 	cs, err := NewCacheService(cfg)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, cs)
// 	cs.Set(t.Context(), "test", "test", 0)
// 	content := ""
// 	cs.Get(t.Context(), "test", &content)
// 	assert.Equal(t, "test", content)

// }
