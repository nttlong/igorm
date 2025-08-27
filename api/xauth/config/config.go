package config

import (
	"os"
	"sync"

	"github.com/spf13/viper"
)

type DbItem struct {
	Dsn string `mapstructure:"dsn"`
}
type Database struct {
	Driver string `mapstructure:"driver"`
	Dsn    string `mapstructure:"dsn"`
}
type Storage struct {
	StorageType string `mapstructure:"type"`
	Location    string `mapstructure:"location"`
}
type Server struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
	Mode string `mapstructure:"mode"`
}
type Config struct {
	Storage Storage `mapstructure:"storage"`
	Server  Server  `mapstructure:"server"`

	Database Database `mapstructure:"database"`
}

var newConfigOnce sync.Once
var config Config

func NewConfig(configPath string) (*Config, error) {

	var err error
	newConfigOnce.Do(func() {

		// viper.SetConfigName("config") // tên file (không có đuôi)
		// viper.SetConfigType("yaml")   // loại file
		if configPath != "" {
			viper.SetConfigFile(configPath)
		} else {
			viper.SetConfigFile("./config.yaml") // đường dẫn tìm file
		}

		viper.AutomaticEnv()
		if err = viper.ReadInConfig(); err != nil {
			return
		}

		if err = viper.Unmarshal(&config); err != nil {
			return
		}
		if os.Getenv("PORT") != "" {
			config.Server.Port = os.Getenv("PORT")
		}

	})
	return &config, err
}
