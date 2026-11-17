package handler

import (
	"strconv"
	"strings"

	"github.com/RomanGulevatiy/skill-badges/internal/badge/domain"
	"github.com/RomanGulevatiy/skill-badges/internal/badge/usecase"
)

const (
	// maxBadgesParamLength limits the raw "b" query parameter length in bytes.
	maxBadgesParamLength = 2048
	// maxBadges limits the number of badges a single request may reference.
	maxBadges = 100
	// maxResponseBytes limits the size of the generated SVG response in bytes.
	maxResponseBytes = 1 << 20

	// defaultPerline is the number of badges per row used when "perline" is invalid.
	defaultPerline = 15
	// minPerline is the smallest accepted "perline" value.
	minPerline = 1
	// maxPerline is the largest accepted "perline" value.
	maxPerline = 32
)

type GetBadgesRequest struct {
	Badges  string
	Perline string
	Theme   string
	Size    string
}

func (r GetBadgesRequest) ToInput() (*usecase.GetBadgesInput, error) {
	if len(r.Badges) > maxBadgesParamLength {
		return nil, domain.ErrBadgesParamTooLong
	}

	rawBadges := strings.Split(r.Badges, ",")

	badges := make([]string, 0, len(rawBadges))
	for _, b := range rawBadges {
		if b = domain.ResolveAlias(b); b != "" {
			badges = append(badges, b)
		}
	}

	if len(badges) > maxBadges {
		return nil, domain.ErrTooManyBadges
	}

	if len(badges) == 0 {
		return nil, domain.ErrNoBadgesProvided
	}

	perline, err := strconv.Atoi(r.Perline)
	if err != nil || perline < minPerline || perline > maxPerline {
		perline = defaultPerline
	}

	return &usecase.GetBadgesInput{
		Badges:  badges,
		Perline: perline,
		Theme:   domain.ValidateTheme(r.Theme),
		Dim:     domain.GetDimensions(r.Size),
	}, nil
}
