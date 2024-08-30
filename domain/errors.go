package domain

import "errors"

var (
	ErrReadingRequestBody = errors.New("failed to read request body")
	ErrParsingRequestBody = errors.New("failed to parse request body")
	ErrNotFound           = errors.New("not found")
)
