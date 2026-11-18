package httpserver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
)

type HTTPServer struct {
	handler http.Handler
	config  Config
	logger  *slog.Logger
}

func New(h http.Handler, cfg Config, log *slog.Logger) *HTTPServer {
	return &HTTPServer{
		handler: h,
		config:  cfg,
		logger:  log,
	}
}

func (s *HTTPServer) Run(ctx context.Context) error {
	ln, err := net.Listen("tcp", s.config.Addr)
	if err != nil {
		return fmt.Errorf("listen HTTP: %w", err)
	}

	return s.serve(ctx, ln)
}

func (s *HTTPServer) serve(ctx context.Context, ln net.Listener) error {
	server := &http.Server{
		Addr:              s.config.Addr,
		Handler:           s.handler,
		ReadHeaderTimeout: s.config.ReadHeaderTimeout,
		ReadTimeout:       s.config.ReadTimeout,
		WriteTimeout:      s.config.WriteTimeout,
		IdleTimeout:       s.config.IdleTimeout,
	}

	serverErr := make(chan error, 1)

	go func() {
		defer close(serverErr)

		s.logger.Info("server started", slog.String("addr", ln.Addr().String()))
		err := server.Serve(ln)
		if errors.Is(err, http.ErrServerClosed) {
			return
		}
		serverErr <- err
	}()

	select {
	case err := <-serverErr:
		if err != nil {
			return fmt.Errorf("serve HTTP: %w", err)
		}
		return nil
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(
		context.WithoutCancel(ctx),
		s.config.ShutDownTimeout,
	)
	defer cancel()

	shutdownErr := server.Shutdown(shutdownCtx)
	if shutdownErr != nil {
		if closeErr := server.Close(); closeErr != nil {
			shutdownErr = errors.Join(shutdownErr, closeErr)
		}
	}

	_ = ln.Close()

	if err := <-serverErr; err != nil {
		if shutdownErr != nil {
			return errors.Join(fmt.Errorf("serve HTTP: %w", err), shutdownErr)
		}
		return fmt.Errorf("serve HTTP: %w", err)
	}

	if shutdownErr != nil {
		return fmt.Errorf("shut down HTTP server: %w", shutdownErr)
	}

	s.logger.Info("server exited gracefully")
	return nil
}
