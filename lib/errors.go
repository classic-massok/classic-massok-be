package lib

import "errors"

// HTTP errors
var (
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbiddden")
	ErrServerError  = errors.New("server error")
)
