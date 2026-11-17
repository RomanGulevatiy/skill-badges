package domain

import "errors"

var (
	// ErrNoBadgesProvided is returned when the "b" parameter contains no badge names.
	ErrNoBadgesProvided = errors.New("at least one badge must be specified in 'b'")
	// ErrBadgeNotFound is returned when a badge is absent from the repository
	// or when none of the requested badges can be resolved.
	ErrBadgeNotFound = errors.New("specified badges were not found in repository")
	// ErrBadgesParamTooLong is returned when the "b" parameter exceeds the allowed length.
	ErrBadgesParamTooLong = errors.New("badge parameter 'b' is too long")
	// ErrTooManyBadges is returned when a request references more badges than allowed.
	ErrTooManyBadges = errors.New("too many badges requested")
)
