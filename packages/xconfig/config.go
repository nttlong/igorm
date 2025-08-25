package xconfig

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

type JWTTokenConfigTTL struct {
	Access  string `yaml:"access"`
	Refresh string `yaml:"refresh"`
}

// Parse string like "120m", "24h", "7d", "2w", "3mo", "1y"
func (ttl *JWTTokenConfigTTL) parseDurationEx(s string) (time.Duration, error) {
	if s == "" {
		return 0, errors.New("empty duration string")
	}

	// tách số và suffix
	numPart := ""
	unitPart := ""
	for i, r := range s {
		if (r >= '0' && r <= '9') || r == '.' {
			numPart += string(r)
		} else {
			unitPart = s[i:]
			break
		}
	}

	if numPart == "" || unitPart == "" {
		return 0, errors.New("invalid duration string: " + s)
	}

	val, err := strconv.ParseFloat(numPart, 64)
	if err != nil {
		return 0, err
	}

	switch strings.ToLower(unitPart) {
	case "ns":
		return time.Duration(val) * time.Nanosecond, nil
	case "us", "µs":
		return time.Duration(val) * time.Microsecond, nil
	case "ms":
		return time.Duration(val) * time.Millisecond, nil
	case "s":
		return time.Duration(val) * time.Second, nil
	case "m":
		return time.Duration(val) * time.Minute, nil
	case "h":
		return time.Duration(val) * time.Hour, nil
	case "d":
		return time.Duration(val*24) * time.Hour, nil
	case "w":
		return time.Duration(val*24*7) * time.Hour, nil
	case "mo", "month":
		// Quy ước 1 tháng = 30 ngày
		return time.Duration(val*24*30) * time.Hour, nil
	case "y", "yr", "year":
		// Quy ước 1 năm = 365 ngày
		return time.Duration(val*24*365) * time.Hour, nil
	default:
		return 0, errors.New("unknown unit: " + unitPart)
	}
}

func (ttl *JWTTokenConfigTTL) GetAccess() (time.Duration, error) {
	var ret time.Duration
	if ttl.Access == "" {
		ret = 15 * time.Minute
		return ret, nil

	}
	return ttl.parseDurationEx(ttl.Access)

}
func (ttl *JWTTokenConfigTTL) GetRefresh() (time.Duration, error) {
	var ret time.Duration
	if ttl.Refresh == "" {
		ret = 7 * 24 * time.Hour
		return ret, nil

	}
	return ttl.parseDurationEx(ttl.Refresh)
}

type JWTTokenConfig struct {
	Secret string            `yaml:"secret"`
	TTL    JWTTokenConfigTTL `yaml:"ttl"`
}
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
	Database Database       `yaml:"database"`
	JwtToken JWTTokenConfig `yaml:"jwt"`
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
