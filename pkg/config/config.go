package config

import (
	"context"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type (
	// Config of service
	Config struct {
		InstanceConfig
		DBConfig
		ClientsConfig
	}

	InstanceConfig struct {
		ServiceName string `env:"SERVICE_NAME,required"`
		HTTPPort    string `env:"HTTP_PORT, default=8084"`
		PublicURL   string `env:"PUBLIC_URL,required"`
		SecretKey   string `env:"SECRET_KEY,required"`
	}

	DBConfig struct {
		PGURL string `env:"PG_URL,required"`
	}

	ClientsConfig struct {
		AIKey string `env:"AI_KEY,required"`
	}
)

// Get
func Get() (Config, error) {
	return parseConfig()
}

func parseConfig() (cfg Config, err error) {
	godotenv.Load()

	err = envconfig.Process(context.Background(), &cfg)
	if err != nil {
		return cfg, fmt.Errorf("fill config: %w", err)
	}

	return cfg, nil
}
