package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCacheService(t *testing.T) {
	cfg, err := NewConfigService("./../config/config.yaml")
	assert.NoError(t, err)
	cs, err := NewCacheService(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, cs)
	cs.Get().Set(t.Context(), "test", "test", 0)
	content := ""
	cs.Get().Get(t.Context(), "test", &content)
	assert.Equal(t, "test", content)

}
