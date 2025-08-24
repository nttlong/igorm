package xconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	config, err := NewConfig(`D:\code\go\news2\igorm\config.yaml`)
	assert.NoError(t, err)
	assert.NotNil(t, config)
}
func BenchmarkTestLoadConfig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := NewConfig(`D:\code\go\news2\igorm\config.yaml`)
		if err != nil {
			b.Error(err)
		}
	}
}
func TestLoadConfigFroRelPath(t *testing.T) {
	config, err := NewConfig(`./../../config.yaml`)
	assert.NoError(t, err)
	assert.NotNil(t, config)
}
func BenchmarkTestLoadConfigFroRelPath(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := NewConfig(`./../../config.yaml`)
		if err != nil {
			b.Error(err)
		}
	}

}
