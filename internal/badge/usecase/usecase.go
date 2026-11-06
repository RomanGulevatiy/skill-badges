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

type Input struct {
	Names      []string
	Perline    int
	Theme      string
	SizePreset string
}

type Output struct {
	SVG string
}

func (u *Usecase) GetBadges(ctx context.Context, in Input) (Output, error) {
	badges := make([]string, 0, len(in.Names))

	theme := domain.ValidateTheme(in.Theme)
	perline := domain.ValidatePerline(in.Perline)

	for _, name := range in.Names {
		canonicalName := domain.ResolveAlias(name)
		if canonicalName == "" {
			continue
		}

		badge, err := u.bRepo.Get(ctx, canonicalName, theme)
		if err != nil {
			u.logger.DebugContext(ctx, "failed to get badge", slog.String("name", canonicalName), slog.Any("error", err))
			continue
		}

		badges = append(badges, badge)
	}

	if len(badges) == 0 {
		return Output{}, domain.ErrBadgeNotFound
	}

	var builder strings.Builder
	builder.Grow(len(badges)*1536 + 200)

	columns := min(len(badges), perline)
	rows := (len(badges) + perline - 1) / perline

	dim := domain.GetDimensions(in.SizePreset)

	width := columns*dim.Size + (columns-1)*dim.Gap
	height := rows*dim.Size + (rows-1)*dim.Gap

	fmt.Fprintf(
		&builder,
		`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" fill="none" viewBox="0 0 %d %d">`,
		width, height, width, height,
	)

	for i, badge := range badges {
		row := i / perline
		col := i % perline

		x := col * (dim.Size + dim.Gap)
		y := row * (dim.Size + dim.Gap)

		fmt.Fprintf(
			&builder,
			`<svg x="%d" y="%d" width="%d" height="%d" viewBox="0 0 256 256">`,
			x, y, dim.Size, dim.Size,
		)
		builder.WriteString(badge)
		builder.WriteString(`</svg>`)
	}

	builder.WriteString(`</svg>`)

	return Output{SVG: builder.String()}, nil
}
