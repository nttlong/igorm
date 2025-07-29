package core_container

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const configPath = "./../../vapi/cmd/config.yaml"

func TestNewContainer(t *testing.T) {
	c := NewCoreContainer(configPath)
	if c == nil {
		t.Error("NewContainer() returned nil")
	}
}
func TestCoreContainer_Config(t *testing.T) {
	c := NewCoreContainer(configPath)
	config := (c.Config.Get()).GetAppConfig()
	assert.Equal(t, "0.0.0.0", config.Host)
}
