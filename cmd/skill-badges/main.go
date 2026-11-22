package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	badgehandler "github.com/RomanGulevatiy/skill-badges/internal/badge/handler"
	badgememory "github.com/RomanGulevatiy/skill-badges/internal/badge/repository/memory"
	badgeusecase "github.com/RomanGulevatiy/skill-badges/internal/badge/usecase"
	"github.com/RomanGulevatiy/skill-badges/internal/config"
	"github.com/RomanGulevatiy/skill-badges/internal/logger"
	"github.com/RomanGulevatiy/skill-badges/internal/middleware"
	httpserver "github.com/RomanGulevatiy/skill-badges/internal/transport/http"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		slog.Error("failed to load config", slog.Any("error", err))
		os.Exit(1)
	}
	log := logger.New(cfg.Logger)

	bRepo, err := badgememory.New("assets")
	if err != nil {
		log.Error("failed to init badge repository", slog.Any("error", err))
		os.Exit(1)
	}

	bUsecase := badgeusecase.New(bRepo)
	bHandler := badgehandler.New(bUsecase, log)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /badges", bHandler.GetBadges)

	handler := middleware.RequestLogger(log)(mux)
	server := httpserver.New(handler, cfg.HTTPServer, log)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err = server.Run(ctx); err != nil {
		log.Error("internal server error", slog.Any("error", err))
		os.Exit(1)
	}
}
