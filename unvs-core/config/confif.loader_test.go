package config

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	LoadConfig("")
	assert.NotEmpty(t, AppConfigInstance)

}
func TestGetCache(t *testing.T) {
	LoadConfig("")
	cache := GetCache()
	assert.NotEmpty(t, cache)
	cache.Set(context.Background(), "test_key", "test_value", 10*time.Second)
	var value string
	cache.Get(context.Background(), "test_key", &value)
	assert.Equal(t, "test_value", value)

}
