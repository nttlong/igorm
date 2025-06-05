package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

// Định nghĩa các struct để chứa cấu hình
// Các tag `mapstructure` giúp Viper ánh xạ đúng các trường từ file YAML
type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
}

type DatabaseConfig struct {
	Driver   string `mapstructure:"driver"` // Thêm trường Driver
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	Timeout  int    `mapstructure:"timeout"`
	SSL      bool   `mapstructure:"ssl"` // Thêm trường SSL
}

type ServerConfig struct {
	Port      int  `mapstructure:"port"`
	DebugMode bool `mapstructure:"debug_mode"`
}

// Struct chính chứa toàn bộ cấu hình ứng dụng
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	Server   ServerConfig   `mapstructure:"server"`
}

// Global variable để lưu trữ cấu hình đã tải
var AppConfigInstance *Config

// LoadConfig tải cấu hình từ file, biến môi trường và giá trị mặc định
func LoadConfig() error {
	// 1. Đặt giá trị mặc định cho các trường cấu hình
	viper.SetDefault("app.name", "DefaultGoApp")
	viper.SetDefault("app.version", "0.0.1")
	viper.SetDefault("database.driver", "mysql") // Mặc định là mysql nếu không có trong config
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306) // Mặc định port cho mysql
	viper.SetDefault("database.user", "root")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.name", "default_db")
	viper.SetDefault("database.timeout", 10)
	viper.SetDefault("database.ssl", false) // Mặc định SSL là false
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.debug_mode", false)

	// 2. Đặt tên file cấu hình (không bao gồm phần mở rộng)
	viper.SetConfigName("config")

	// 3. Đặt loại file cấu hình
	viper.SetConfigType("yaml")

	// 4. Đặt đường dẫn tìm kiếm file cấu hình
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../..")
	viper.AddConfigPath("/etc/appname")
	viper.AddConfigPath("$HOME/.appname")

	// 5. Cấu hình để đọc biến môi trường
	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 6. Đọc cấu hình
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("Cảnh báo: Không tìm thấy file cấu hình. Sử dụng giá trị mặc định và biến môi trường. Lỗi: %v", err)
		} else {
			return fmt.Errorf("lỗi khi đọc file cấu hình: %w", err)
		}
	}

	// 7. Binding cấu hình vào struct
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return fmt.Errorf("không thể binding cấu hình vào struct: %w", err)
	}

	// Gán cấu hình đã tải vào biến global để dễ dàng truy cập từ các nơi khác
	AppConfigInstance = &cfg

	log.Println("Cấu hình đã được tải thành công.")
	return nil
}
