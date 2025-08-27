package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	cfg, err := NewConfig("")
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
}
