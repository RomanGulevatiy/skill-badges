package handler

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/RomanGulevatiy/skill-badges/internal/badge/domain"
	"github.com/RomanGulevatiy/skill-badges/internal/badge/usecase"
)

const maxBadgeParamLen = 1024

type Usecase interface {
	GetBadges(ctx context.Context, in usecase.Input) (usecase.Output, error)
}

type Handler struct {
	bUsecase Usecase
	logger   *slog.Logger
}

func New(bUsecase Usecase, log *slog.Logger) *Handler {
	return &Handler{bUsecase: bUsecase, logger: log}
}

func (h *Handler) GetBadges(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	b := q.Get("b")
	if b == "" {
		http.Error(w, domain.ErrMissingBadgeParam.Error(), http.StatusBadRequest)
		return
	}
	if len(b) > maxBadgeParamLen {
		http.Error(w, domain.ErrBadgeParamTooLong.Error(), http.StatusBadRequest)
		return
	}

	names := strings.Split(b, ",")
	theme := q.Get("theme")
	perline, _ := strconv.Atoi(q.Get("perline"))
	size := q.Get("size")

	input := usecase.Input{
		Names:      names,
		Theme:      theme,
		Perline:    perline,
		SizePreset: size,
	}

	output, err := h.bUsecase.GetBadges(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrBadgeNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			h.logger.ErrorContext(r.Context(), "get badges failed", slog.Any("error", err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'unsafe-inline'")
	if _, err := io.WriteString(w, output.SVG); err != nil {
		h.logger.WarnContext(r.Context(), "write response failed", slog.Any("error", err))
	}
}
