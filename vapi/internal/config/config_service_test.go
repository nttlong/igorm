package config

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigService(t *testing.T) {
	cfg, err := NewConfigService("./../config/config.yaml")
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "./tmp/badger", cfg.config.Badger.Path)
	absPath, err := filepath.Abs(cfg.config.Badger.Path)
	if err != nil {
		panic(err)
	}
	fmt.Println("Absolute path:", absPath)

}
func TestCacheService_DI(t *testing.T) {
	// root := vdi.NewRootContainer()

}
