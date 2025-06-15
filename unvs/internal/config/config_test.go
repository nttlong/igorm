package config_test

import (
	"fmt"
	"log"
	"testing"
	"unvs/internal/config"
)

func TestConfig(t *testing.T) {

	if err := config.LoadConfig(); err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	fmt.Println("--- loaded config ---")
	fmt.Printf("App Name: %s\n", config.AppConfigInstance.App.Name)
	fmt.Printf("DB Driver: %s\n", config.AppConfigInstance.Database.Driver)
	fmt.Printf("DB Host: %s:%d\n", config.AppConfigInstance.Database.Host, config.AppConfigInstance.Database.Port)
	fmt.Printf("DB SSL: %t\n", config.AppConfigInstance.Database.SSL)
	fmt.Printf("Server Port: %s\n", config.AppConfigInstance.Server.Port)
}
