package config

// Định nghĩa các struct để chứa cấu hình
// Các tag `mapstructure` giúp Viper ánh xạ đúng các trường từ file YAML
type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
}

type DatabaseConfig struct {
	Driver         string `mapstructure:"driver"` // Thêm trường Driver
	Host           string `mapstructure:"host"`
	Port           int    `mapstructure:"port"`
	User           string `mapstructure:"user"`
	Password       string `mapstructure:"password"`
	Name           string `mapstructure:"name"`
	Timeout        int    `mapstructure:"timeout"`
	SSL            bool   `mapstructure:"ssl"` // Thêm trường SSL
	IsMultiTenancy bool   `mapstructure:"isMultiTenancy"`
}

type ServerConfig struct {
	Port      string `mapstructure:"port"`
	Bind      string `mapstructure:"bind"`
	DebugMode bool   `mapstructure:"debug_mode"`
}

// Struct chính chứa toàn bộ cấu hình ứng dụng
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	Server   ServerConfig   `mapstructure:"server"`
	// cache type string value is bagger, in-memory, redis, memcached
	Cache           string     `mapstructure:"cache"`
	MemcachedServer []string   `mapstructure:"memcached"`
	Redis           RedisCache `mapstructure:"redis"`
	Logs            string     `mapstructure:"logs"`
	EncryptionKey   string     `mapstructure:"encryptionKey"`
	IsDebug         bool       `mapstructure:"isDebug"`
}
type RedisCache struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	Timeout  int    `mapstructure:"timeout"`
}

// Global variable để lưu trữ cấu hình đã tải
var AppConfigInstance *Config
