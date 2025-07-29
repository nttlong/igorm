package config_service

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const configPath = "./../../vapi/cmd/config.yaml"

func TestConfigService(t *testing.T) {
	cfg, err := NewConfigService(configPath)
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.Equal(t, "./tmp/badger", cfg.GetAppConfig().Badger.Path)
	absPath, err := filepath.Abs(cfg.GetAppConfig().Badger.Path)
	if err != nil {
		panic(err)
	}
	fmt.Println("Absolute path:", absPath)

}
func TestCacheService_DI(t *testing.T) {
	// root := vdi.NewRootContainer()

}
