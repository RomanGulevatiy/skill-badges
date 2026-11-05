package main

import (
	"log/slog"
	"net/http"
	"os"

	badgehandler "github.com/RomanGulevatiy/skill-badges/internal/badge/handler"
	badgememory "github.com/RomanGulevatiy/skill-badges/internal/badge/repository/memory"
	badgeusecase "github.com/RomanGulevatiy/skill-badges/internal/badge/usecase"
	"github.com/RomanGulevatiy/skill-badges/internal/config"
	"github.com/RomanGulevatiy/skill-badges/internal/logger"
	"github.com/RomanGulevatiy/skill-badges/internal/middleware"
	httpserver "github.com/RomanGulevatiy/skill-badges/internal/transport/http"
)

func main() {
	cfg := config.MustNew()
	log := logger.New(cfg.Logger)

	bRepo, err := badgememory.New("assets", log)
	if err != nil {
		log.Error("failed to init badge repository", slog.Any("error", err))
		os.Exit(1)
	}

	bUsecase := badgeusecase.New(bRepo, log)
	bHandler := badgehandler.New(bUsecase, log)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/badges", bHandler.GetBadges)

	logMiddleware := middleware.RequestLogger(log)
	h := logMiddleware(mux)

	server := httpserver.New(h, cfg.HTTPServer, log)

	if err := server.Run(); err != nil {
		log.Error("internal server error", slog.Any("error", err))
		os.Exit(1)
	}
}
