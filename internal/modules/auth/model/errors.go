package model

import (
	"errors"
)

var (
	ErrAccountNotFound   = errors.New("account not found")
	ErrEmailAlreadyTaken = errors.New("email already taken")

	ErrSessionNotFound    = errors.New("session not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrSessionExpired     = errors.New("session expired")
)
