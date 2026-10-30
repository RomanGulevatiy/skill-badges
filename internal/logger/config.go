package logger

import (
	"log/slog"
)

type Config struct {
	Level slog.Level `env:"LOGGER_LEVEL" envDefault:"INFO"`
}
