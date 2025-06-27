package config

import (
	"fmt"
	"log/slog"

	http_server "github.com/Util787/task-manager/pkg/http-server"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Env        string `env:"ENV" envDefault:"prod"`
	HttpServer http_server.Config
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	slog.Info("Configuration loaded successfully", // for debug, delete later
		"env", cfg.Env,
		"http_port", cfg.HttpServer.Port,
		"http_read_header_timeout", cfg.HttpServer.ReadHeaderTimeout,
		"http_write_timeout", cfg.HttpServer.WriteTimeout,
		"http_read_timeout", cfg.HttpServer.ReadTimeout,
	)

	return cfg, nil
}
