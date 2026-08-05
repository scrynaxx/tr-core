package security

import "errors"

var (
	ErrInvalidToken   = errors.New("invalid token")
	ErrSessionRevoked = errors.New("session revoked")
	ErrForbidden      = errors.New("forbidden")
)
