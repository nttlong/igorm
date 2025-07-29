package config_service

import (
	"time"
)

type LogConfig struct {
	FileName   string `mapstructure:"fileName"`
	MaxSize    int    `mapstructure:"maxSize"`
	MaxAge     int    `mapstructure:"maxAge"`
	MaxBackups int    `mapstructure:"maxBackups"`
	Compress   bool   `mapstructure:"compress"`
	AlsoStdout bool   `mapstructure:"alsoStdout"`
}
type Config struct {
	CacheType  string          `mapstructure:"cacheType"`
	DriverType string          `mapstructure:"driverType"`
	Redis      RedisConfig     `mapstructure:"redis"`
	Memcached  MemcachedConfig `mapstructure:"memcached"`
	Badger     BadgerConfig    `mapstructure:"badger"`
	InMemory   InMemoryConfig  `mapstructure:"inmemory"`
	Database   DatabaseConfig  `mapstructure:"database"`
	Host       string          `mapstructure:"host"`
	Port       int             `mapstructure:"port"`
}

type RedisConfig struct {
	Nodes     string        `mapstructure:"nodes"`
	Password  string        `mapstructure:"password"`
	PrefixKey string        `mapstructure:"prefixKey"`
	Timeout   time.Duration `mapstructure:"timeout"`
	DB        int           `mapstructure:"db"`
}

type MemcachedConfig struct {
	Nodes     string `mapstructure:"nodes"`
	PrefixKey string `mapstructure:"prefixKey"`
}

type BadgerConfig struct {
	Path      string `mapstructure:"path"`
	PrefixKey string `mapstructure:"prefixKey"`
}

type InMemoryConfig struct {
	DefaultTTL      time.Duration `mapstructure:"defaultTTL"`
	CleanupInterval time.Duration `mapstructure:"cleanupInterval"`
}

type DatabaseConfig struct {
	Driver string `mapstructure:"driver"`
	DSN    string `mapstructure:"dsn"`
}

// type ConfigService interface {
// 	GetConfig() *Config
// }

// type configService struct {
// 	cfg *Config
// }

// func NewConfigService(configPath string) (ConfigService, error) {
// 	v := viper.New()
// 	v.SetConfigFile(configPath)
// 	v.SetConfigType("yaml")

// 	if err := v.ReadInConfig(); err != nil {
// 		return nil, fmt.Errorf("cannot read config file: %w", err)
// 	}

// 	var cfg Config
// 	if err := v.Unmarshal(&cfg); err != nil {
// 		return nil, fmt.Errorf("cannot unmarshal config: %w", err)
// 	}

// 	return &configService{cfg: &cfg}, nil
// }

// func (s *configService) GetConfig() *Config {
// 	return s.cfg
// }
