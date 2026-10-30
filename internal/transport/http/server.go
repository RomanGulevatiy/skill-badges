package httpserver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
)

type HTTPServer struct {
	mux    *http.ServeMux
	config Config
	logger *slog.Logger
}

func New(mux *http.ServeMux, cfg Config, log *slog.Logger) *HTTPServer {
	return &HTTPServer{
		mux:    mux,
		config: cfg,
		logger: log,
	}
}

func (s *HTTPServer) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	server := &http.Server{
		Addr:              s.config.Addr,
		Handler:           s.mux,
		ReadHeaderTimeout: s.config.ReadHeaderTimeout,
		ReadTimeout:       s.config.ReadTimeout,
		WriteTimeout:      s.config.WriteTimeout,
		IdleTimeout:       s.config.IdleTimeout,
	}

	serverErr := make(chan error, 1)

	go func() {
		defer close(serverErr)

		s.logger.Info("server started", slog.String("addr", s.config.Addr))
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		if err != nil {
			return fmt.Errorf("listen and serve HTTP: %w", err)
		}
	case <-ctx.Done():
		s.logger.Info("context cancelled")
	}

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		s.config.ShutDownTimeout,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		if closeErr := server.Close(); closeErr != nil {
			return errors.Join(err, closeErr)
		}
		return fmt.Errorf("shut down HTTP server: %w", err)
	}

	<-serverErr

	s.logger.Info("server exited gracefully")
	return nil
}
