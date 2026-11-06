package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/RomanGulevatiy/skill-badges/internal/badge/domain"
)

type Repository interface {
	Get(ctx context.Context, name, theme string) (string, error)
}

type Usecase struct {
	bRepo  Repository
	logger *slog.Logger
}

func New(bRepo Repository, log *slog.Logger) *Usecase {
	return &Usecase{
		bRepo:  bRepo,
		logger: log,
	}
}

type GetBadgesInput struct {
	Badges  []string
	Perline int
	Theme   string
	Dim     domain.Dimensions
}

func (u *Usecase) GetBadges(ctx context.Context, in GetBadgesInput) (string, error) {
	badges := make([]string, 0, len(in.Badges))

	for _, badge := range in.Badges {
		b, err := u.bRepo.Get(ctx, badge, in.Theme)
		if err != nil {
			u.logger.DebugContext(ctx, "failed to get badge", slog.String("badge", badge), slog.Any("error", err))
			continue
		}
		badges = append(badges, b)
	}

	if len(badges) == 0 {
		return "", domain.ErrBadgeNotFound
	}

	return render(badges, in.Perline, in.Dim), nil
}

func render(badges []string, perline int, dim domain.Dimensions) string {
	columns := min(len(badges), perline)
	rows := (len(badges) + perline - 1) / perline

	width := columns*dim.Size + (columns-1)*dim.Gap
	height := rows*dim.Size + (rows-1)*dim.Gap

	capacity := 256
	for _, badge := range badges {
		capacity += len(badge) + 90
	}

	var b strings.Builder
	b.Grow(capacity)

	fmt.Fprintf(&b,
		`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" fill="none" viewBox="0 0 %d %d">`,
		width, height, width, height,
	)

	for i, badge := range badges {
		x := (i % perline) * (dim.Size + dim.Gap)
		y := (i / perline) * (dim.Size + dim.Gap)

		fmt.Fprintf(&b,
			`<svg x="%d" y="%d" width="%d" height="%d" viewBox="0 0 256 256">`,
			x, y, dim.Size, dim.Size,
		)
		b.WriteString(badge)
		b.WriteString(`</svg>`)
	}

	b.WriteString(`</svg>`)
	return b.String()
}
