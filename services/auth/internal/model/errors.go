package model

import "errors"

var (
	ErrSessionNotFound    = errors.New("session not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrSessionExpired     = errors.New("session expired")
)
