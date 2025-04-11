package jwt

import "errors"

var (
	ErrInvalidToken = errors.New("token format is invalid. Expected: Bearer <token>")
)
