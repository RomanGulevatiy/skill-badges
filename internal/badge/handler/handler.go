package handler

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/RomanGulevatiy/skill-badges/internal/badge/domain"
	"github.com/RomanGulevatiy/skill-badges/internal/badge/usecase"
)

type Usecase interface {
	GetBadges(ctx context.Context, in usecase.GetBadgesInput) (string, error)
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

	req := GetBadgesRequest{
		Badges:  q.Get("b"),
		Perline: q.Get("perline"),
		Theme:   q.Get("theme"),
		Size:    q.Get("size"),
	}

	in, err := req.Validate()
	if err != nil {
		h.logger.DebugContext(r.Context(), "invalid query parameters", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	out, err := h.bUsecase.GetBadges(r.Context(), *in)
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
	if _, err = io.WriteString(w, out); err != nil {
		h.logger.WarnContext(r.Context(), "write response failed", slog.Any("error", err))
	}
}
