package handler

import (
	"strconv"
	"strings"

	"github.com/RomanGulevatiy/skill-badges/internal/badge/domain"
	"github.com/RomanGulevatiy/skill-badges/internal/badge/usecase"
)

type GetBadgesRequest struct {
	Badges  string
	Perline string
	Theme   string
	Size    string
}

func (r GetBadgesRequest) Validate() (*usecase.GetBadgesInput, error) {
	rawBadges := strings.Split(r.Badges, ",")

	badges := make([]string, 0, len(rawBadges))
	for _, b := range rawBadges {
		if b = domain.ResolveAlias(b); b != "" {
			badges = append(badges, b)
		}
	}

	switch {
	case len(badges) == 0:
		return nil, domain.ErrNoBadgesProvided
	case len(badges) > 2048:
		return nil, domain.ErrBadgeParamTooLong
	}

	perline, err := strconv.Atoi(r.Perline)
	if err != nil || perline < 1 || perline > 32 {
		perline = 15
	}

	return &usecase.GetBadgesInput{
		Badges:  badges,
		Perline: perline,
		Theme:   domain.ValidateTheme(r.Theme),
		Dim:     domain.GetDimensions(r.Size),
	}, nil
}
