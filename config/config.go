package config

import (
	"context"

	"github.com/kimnguyenlong/ketoz/pkg/keto"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	Service ServiceConfig `env:",prefix=SERVICE_"`
	Keto    keto.Config   `env:",prefix=KETO_"`
}

type ServiceConfig struct {
	Host     string `env:"HOST, default=0.0.0.0"`
	Port     int    `env:"PORT, default=8000"`
	LogLevel string `env:"LOG_LEVEL, default=INFO"`
}

func Load() (*Config, error) {
	ctx := context.Background()

	var c Config
	if err := envconfig.Process(ctx, &c); err != nil {
		return nil, err
	}
	return &c, nil
}
