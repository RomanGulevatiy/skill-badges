package httpserver

import "time"

type Config struct {
	Addr              string
	ShutDownTimeout   time.Duration
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
}

// TODO: Apply environment variables instead of default values
func NewConfig() Config {
	return Config{
		Addr:              ":8080",
		ShutDownTimeout:   15 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}
