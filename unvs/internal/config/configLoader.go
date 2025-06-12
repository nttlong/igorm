package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

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
	viper.SetDefault("cache", "bagger")
	viper.SetDefault("logs", "./logs/app.log")

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
			log.Printf("Warning: can not find config file. Use default value and environment variable instead.")

		} else {
			return fmt.Errorf("read config file error: %w", err)
		}
	}

	// 7. Binding cấu hình vào struct
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return fmt.Errorf("can not bind config to struct: %w", err)
	}

	// Gán cấu hình đã tải vào biến global để dễ dàng truy cập từ các nơi khác
	AppConfigInstance = &cfg

	log.Println("config loaded successfully.")
	return nil
}
