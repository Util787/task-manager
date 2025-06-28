package config

import (
	"fmt"

	http_server "github.com/Util787/task-manager/pkg/http-server"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Env           string `env:"ENV" envDefault:"prod"`
	HttpServerCfg http_server.Config
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config with err: %w, using default values", err)
	}

	if cfg.Env != "prod" && cfg.Env != "dev" && cfg.Env != "local" {
		return nil, fmt.Errorf("invalid environment: %s, must be prod, dev or local", cfg.Env)
	}

	return cfg, nil
}
