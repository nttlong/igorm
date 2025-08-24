package xconfig

import (
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

type DbItem struct {
	Dsn     string `yaml:"dsn"`
	Manager string `yaml:"manager"`
}

// Database config . when app start it will read this config from yaml file
type Database struct {
	Postgres DbItem `yaml:"postgres"`
	Driver   string `yaml:"driver"`
}
type Config struct {
	Database Database `yaml:"database"`
	Manager  string   `yaml:"manager"`
}
type initNewConfig struct {
	val  *Config
	err  error
	once sync.Once
}

var cacheNewConfig sync.Map

func NewConfig(configFile string) (*Config, error) {
	absFilePath, err := filepath.Abs(configFile)
	if err != nil {
		return nil, err
	}
	actualConfig, _ := cacheNewConfig.LoadOrStore(absFilePath, &initNewConfig{})
	load := actualConfig.(*initNewConfig)
	load.once.Do(func() {
		config := &Config{}
		load.val = config

		if err != nil {
			load.err = err
			return
		}
		// usi ioutil.ReadFile to read the file
		data, err := os.ReadFile(absFilePath)
		if err != nil {
			load.err = err
			return
		}
		yaml.Unmarshal(data, config)
		load.val = config
		return
	})
	return load.val, load.err

}
