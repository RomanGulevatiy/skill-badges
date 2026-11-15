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

func (r GetBadgesRequest) ToInput() (*usecase.GetBadgesInput, error) {
	rawBadges := strings.Split(r.Badges, ",")

	badges := make([]string, 0, len(rawBadges))
	for _, b := range rawBadges {
		if b = domain.ResolveAlias(b); b != "" {
			badges = append(badges, b)
		}
	}

	if len(badges) == 0 {
		return nil, domain.ErrNoBadgesProvided
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
