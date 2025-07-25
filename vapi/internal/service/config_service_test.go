package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigService(t *testing.T) {
	cfg, err := NewConfigService("./../config/config.yaml")
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "./tmp/badger", cfg.config.Badger.Path)

}
