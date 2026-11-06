package domain

import "errors"

var (
	ErrNoBadgesProvided  = errors.New("at least one badge must be specified in 'b'")
	ErrBadgeNotFound     = errors.New("specified badges were not found in repository")
	ErrBadgeParamTooLong = errors.New("badge parameter is too long")
)
