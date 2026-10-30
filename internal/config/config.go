package config

import (
	"fmt"

	"github.com/RomanGulevatiy/skill-badges/internal/logger"
	httpserver "github.com/RomanGulevatiy/skill-badges/internal/transport/http"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Logger     logger.Config
	HTTPServer httpserver.Config
}

func New() (*Config, error) {
	_ = godotenv.Load()

	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}

func MustNew() *Config {
	config, err := New()

	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	return config
}
