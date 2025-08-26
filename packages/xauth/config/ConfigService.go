package config

import (
	"xconfig"
)

type ConfigService interface {
	Get() *xconfig.Config
}
type YamlConfigService struct {
	data *xconfig.Config
}

func (configService *YamlConfigService) New() (ConfigService, error) {
	cfg, err := xconfig.NewConfig("./../config.yaml")
	if err != nil {
		return nil, err
	}
	configService.data = cfg
	return configService, nil
}
func (configService *YamlConfigService) Get() *xconfig.Config {
	return configService.data
}
