package domain

import "errors"

var (
	ErrMissingBadgesParam = errors.New("badge parameter is required")
	ErrBadgeNotFound      = errors.New("badge not found")
)
