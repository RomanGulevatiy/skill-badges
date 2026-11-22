package usecase

import (
	"context"
	"strconv"
	"strings"

	"github.com/RomanGulevatiy/skill-badges/internal/badge/domain"
)

type Repository interface {
	Get(ctx context.Context, name, theme string) (string, error)
}

type Usecase struct {
	bRepo Repository
}

func New(bRepo Repository) *Usecase {
	return &Usecase{bRepo: bRepo}
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

	var scratch [160]byte

	buf := append(scratch[:0], `<svg xmlns="http://www.w3.org/2000/svg" width="`...)
	buf = strconv.AppendInt(buf, int64(width), 10)
	buf = append(buf, `" height="`...)
	buf = strconv.AppendInt(buf, int64(height), 10)
	buf = append(buf, `" fill="none" viewBox="0 0 `...)
	buf = strconv.AppendInt(buf, int64(width), 10)
	buf = append(buf, ' ')
	buf = strconv.AppendInt(buf, int64(height), 10)
	buf = append(buf, `">`...)
	b.Write(buf)

	for i, badge := range badges {
		x := (i % perline) * (dim.Size + dim.Gap)
		y := (i / perline) * (dim.Size + dim.Gap)

		buf = append(scratch[:0], `<svg x="`...)
		buf = strconv.AppendInt(buf, int64(x), 10)
		buf = append(buf, `" y="`...)
		buf = strconv.AppendInt(buf, int64(y), 10)
		buf = append(buf, `" width="`...)
		buf = strconv.AppendInt(buf, int64(dim.Size), 10)
		buf = append(buf, `" height="`...)
		buf = strconv.AppendInt(buf, int64(dim.Size), 10)
		buf = append(buf, `" viewBox="0 0 256 256">`...)
		b.Write(buf)
		b.WriteString(badge)
		b.WriteString(`</svg>`)
	}

	b.WriteString(`</svg>`)
	return b.String()
}
