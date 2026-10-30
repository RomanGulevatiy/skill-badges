package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/RomanGulevatiy/skill-badges/internal/config"
	"github.com/RomanGulevatiy/skill-badges/internal/handler"
	"github.com/RomanGulevatiy/skill-badges/internal/logger"
	"github.com/RomanGulevatiy/skill-badges/internal/repository/inmemory"
	httpserver "github.com/RomanGulevatiy/skill-badges/internal/transport/http"
	"github.com/RomanGulevatiy/skill-badges/internal/usecase"
)

func main() {
	config := config.MustNew()

	logger := logger.New(config.Logger)

	bInmemoryRepo, err := inmemory.NewBadgeInmemoryRepo("assets", logger)
	if err != nil {
		logger.Error("failed to init badge in-memory repo", slog.Any("error", err))
		os.Exit(1)
	}

	bUsecase := usecase.NewBadgeUsecase(bInmemoryRepo, logger)
	bHandler := handler.NewBadgeHandler(bUsecase, logger)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/badges", bHandler.GetBadges)

	server := httpserver.New(mux, config.HTTPServer, logger)

	if err := server.Run(); err != nil {
		logger.Error("internal server error", slog.Any("error", err))
		os.Exit(1)
	}
}
