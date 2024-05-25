package utils

import "errors"

var (
	ErrUserIDRequired        = errors.New("user id required")
	ErrProjectIDRequired     = errors.New("project id required")
	ErrEmailRequired         = errors.New("email required")
	ErrMethodNotAllowed      = errors.New("invalid request method")
	ErrInvalidRequestPayload = errors.New("invalid request payload")
	ErrReadingResponse       = errors.New("error reading response BODY")
	ErrCreatingSession       = errors.New("error creating session")
	ErrUnauthorized          = errors.New("invalid credentials")
)
