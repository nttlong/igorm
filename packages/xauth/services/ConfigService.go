package services

import (
	"xconfig"
)

type ConfigService struct {
	Data *xconfig.Config
}

func (configService *ConfigService) New() error {
	cfg, err := xconfig.NewConfig("./../../../config.yaml")
	if err != nil {
		return err
	}
	configService.Data = cfg
	return nil
}
