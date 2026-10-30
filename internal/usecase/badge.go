package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/RomanGulevatiy/skill-badges/internal/domain"
)

type BadgeInmemoryRepo interface {
	GetBadge(ctx context.Context, name string) (string, error)
}

type BadgeUsecase struct {
	bInmemoryRepo BadgeInmemoryRepo
	logger        *slog.Logger
}

func NewBadgeUsecase(bInmemoryRepo BadgeInmemoryRepo, log *slog.Logger) *BadgeUsecase {
	return &BadgeUsecase{
		bInmemoryRepo: bInmemoryRepo,
		logger:        log,
	}
}

type (
	Input struct {
		BadgeParam   string
		PerlineParam int
		GapParam     int
		SizeParam    int
	}

	Output string
)

func (u *BadgeUsecase) GetBadges(ctx context.Context, in Input) (Output, error) {
	names := strings.Split(in.BadgeParam, ",")
	badges := make([]string, 0, len(names))

	for _, name := range names {
		trimmed := strings.TrimSpace(name)
		if trimmed == "" {
			continue
		}

		badge, err := u.bInmemoryRepo.GetBadge(ctx, trimmed)
		if err != nil {
			u.logger.DebugContext(ctx, "failed to get badge", slog.String("name", trimmed), slog.Any("error", err))
			continue
		}

		badges = append(badges, badge)
	}

	if len(badges) == 0 {
		return "", domain.ErrBadgeNotFound
	}

	var builder strings.Builder
	builder.Grow(len(badges)*1536 + 200)

	columns := min(len(badges), in.PerlineParam)
	rows := (len(badges) + in.PerlineParam - 1) / in.PerlineParam

	width := columns*in.SizeParam + (columns-1)*in.GapParam
	height := rows*in.SizeParam + (rows-1)*in.GapParam

	builder.WriteString(fmt.Sprintf(
		`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" fill="none" viewBox="0 0 %d %d">`,
		width, height, width, height,
	))

	for i, badge := range badges {
		row := i / in.PerlineParam
		col := i % in.PerlineParam

		x := col * (in.SizeParam + in.GapParam)
		y := row * (in.SizeParam + in.GapParam)

		builder.WriteString(fmt.Sprintf(
			`<svg x="%d" y="%d" width="%d" height="%d" viewBox="0 0 256 256">`,
			x, y, in.SizeParam, in.SizeParam,
		))
		builder.WriteString(badge)
		builder.WriteString(`</svg>`)
	}

	builder.WriteString(`</svg>`)

	return Output(builder.String()), nil
}
