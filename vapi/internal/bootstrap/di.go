package bootstrap

import (
	"log"

	"vapi/internal/config"
	"vapi/internal/service"
	"vdi"
)

func RegisterServices(root vdi.RootContainer) { // <-- KHÔNG dùng *vdi.Container
	// Load config từ file YAML
	
	root.RegisterSingleton(func() *service.ConfigService {
		return &Logger{ID: "singleton"}
	})

	cfg, err := service.NewConfigService("./../config/config.yaml")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	// Đăng ký ConfigService dưới dạng singleton (instance đã khởi tạo)
	root.RegisterInstance(cfg) // root dang la vdi.Container va
	// vdi.Container duoc dinh nghia nhu sau
	// type Container interface {
	//ResolveByType(t reflect.Type) (any, error)
	//}, thi ro rang la bi loi

	// Đăng ký CacheService
	root.RegisterSingleton(func(cfg *config.ConfigService) (*cache.CacheService, error) {
		return cache.NewCacheService(cfg)
	})
}
