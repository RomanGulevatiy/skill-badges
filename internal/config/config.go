package config

import (
	"errors"
	"fmt"
	"os"

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
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("failed to load .env: %w", err)
	}

	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return &cfg, nil
}
