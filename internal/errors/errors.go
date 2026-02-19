package apperrors

import "errors"

var (
	ErrNotFound = errors.New("resource not found")
	ErrBadID    = errors.New("invalid id")
)
