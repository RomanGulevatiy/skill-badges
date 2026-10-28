package handler

import (
	"context"
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
}

func NewBadgeHandler(bUsecase BadgeUsecase) *BadgeHandler {
	return &BadgeHandler{bUsecase: bUsecase}
}

func (h *BadgeHandler) GetBadges(w http.ResponseWriter, r *http.Request) {
	queryParam := r.URL.Query()

	b := queryParam.Get("b")
	if b == "" {
		http.Error(w, domain.ErrMissingBadgesParam.Error(), http.StatusBadRequest)
		return
	}

	perline, err := strconv.Atoi(queryParam.Get("perline"))
	if err != nil || perline < 1 || perline > 32 {
		perline = 15 // Default variable for perline
	}

	gap, err := strconv.Atoi(queryParam.Get("gap"))
	if err != nil || gap < 0 || gap > 100 {
		gap = 8 // Default variable for gap
	}

	size, err := strconv.Atoi(queryParam.Get("size"))
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(output))
}
