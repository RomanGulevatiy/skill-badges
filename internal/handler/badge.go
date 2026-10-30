package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/RomanGulevatiy/skill-badges/internal/domain"
	"github.com/RomanGulevatiy/skill-badges/internal/usecase"
)

type BadgeUsecase interface {
	GetBadges(ctx context.Context, in usecase.Input) (usecase.Output, error)
}

type BadgeHandler struct {
	bUsecase BadgeUsecase
	logger   *slog.Logger
}

func NewBadgeHandler(bUsecase BadgeUsecase, log *slog.Logger) *BadgeHandler {
	return &BadgeHandler{bUsecase: bUsecase, logger: log}
}

func (h *BadgeHandler) GetBadges(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	b := q.Get("b")
	if b == "" {
		http.Error(w, domain.ErrMissingBadgesParam.Error(), http.StatusBadRequest)
		return
	}

	perline, err := strconv.Atoi(q.Get("perline"))
	if err != nil || perline < 1 || perline > 32 {
		perline = 15 // Default variable for perline
	}

	gap, err := strconv.Atoi(q.Get("gap"))
	if err != nil || gap < 0 || gap > 100 {
		gap = 8 // Default variable for gap
	}

	size, err := strconv.Atoi(q.Get("size"))
	if err != nil || size < 1 || size > 256 {
		size = 48 // Default variable for size
	}

	input := usecase.Input{
		BadgeParam:   b,
		PerlineParam: perline,
		GapParam:     gap,
		SizeParam:    size,
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
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(output))
}
