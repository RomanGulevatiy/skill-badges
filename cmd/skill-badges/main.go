package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/RomanGulevatiy/skill-badges/internal/handler"
	"github.com/RomanGulevatiy/skill-badges/internal/repository/inmemory"
	httpserver "github.com/RomanGulevatiy/skill-badges/internal/transport/http"
	"github.com/RomanGulevatiy/skill-badges/internal/usecase"
)

func main() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelWarn,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))

	bInmemoryRepo, err := inmemory.NewBadgeInmemoryRepo("assets", logger)
	if err != nil {
		logger.Error("failed to init badge in-memory repo", slog.Any("error", err))
		os.Exit(1)
	}

	bUsecase := usecase.NewBadgeUsecase(bInmemoryRepo, logger)
	bHandler := handler.NewBadgeHandler(bUsecase)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/badges", bHandler.GetBadges)

	srvConfig := httpserver.NewConfig()
	server := httpserver.New(mux, srvConfig, logger)

	if err := server.Run(); err != nil {
		logger.Error("internal server error", slog.Any("error", err))
		os.Exit(1)
	}
}
