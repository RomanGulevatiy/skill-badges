package domain

import "errors"

var (
	ErrMissingBadgeParam = errors.New("badge parameter is required")
	ErrBadgeNotFound     = errors.New("badge not found")
)
