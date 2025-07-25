package service

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Redis struct {
	Nodes     string        `mapstructure:"nodes"`
	Password  string        `mapstructure:"password"`
	PrefixKey string        `mapstructure:"prefixKey"`
	DB        int           `mapstructure:"db"`
	Timeout   time.Duration `mapstructure:"timeout"`
}

type Memcached struct {
	Nodes     string `mapstructure:"nodes"`
	PrefixKey string `mapstructure:"prefixKey"`
}

type Badger struct {
	Path      string `mapstructure:"path"`
	PrefixKey string `mapstructure:"prefixKey"`
}

type InMemory struct {
	DefaultTTL      time.Duration `mapstructure:"defaultTTL"`
	CleanupInterval time.Duration `mapstructure:"cleanupInterval"`
}

type AppConfig struct {
	CacheType string    `mapstructure:"cacheType"` // redis | memcached | badger | inmemory
	Redis     Redis     `mapstructure:"redis"`
	Memcached Memcached `mapstructure:"memcached"`
	Badger    Badger    `mapstructure:"badger"`
	InMemory  InMemory  `mapstructure:"inmemory"`

	// Future: Add database config here
}

type ConfigService struct {
	config *AppConfig
}

func NewConfigService(configFilePath string) (*ConfigService, error) {
	v := viper.New()
	v.SetConfigFile(configFilePath)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var c AppConfig
	if err := v.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &ConfigService{config: &c}, nil
}

func (s *ConfigService) Get() *AppConfig {
	return s.config
}
