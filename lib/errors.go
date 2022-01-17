package lib

import "errors"

// HTTP errors
var (
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrServerError  = errors.New("internal server error")
)
