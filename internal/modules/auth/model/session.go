package model

import (
	"errors"
	"time"

	"uuid"
)

type Session struct {
	ID          uuid.UUID `json:"id"`
	AccountID   uuid.UUID `json:"account_id"`
	UserAgent   string    `json:"user_agent"`
	RefreshHash string    `json:"refresh_hash"`
	ExpiresAt   time.Time `json:"expires_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewSession(accountID uuid.UUID, userAgent, refreshHash string, lifetime time.Duration) (Session, error) {
	if accountID == uuid.Nil() {
		return Session{}, errors.New("employee ID is empty")
	}
	if userAgent == "" {
		return Session{}, errors.New("user agent is empty")
	}
	if refreshHash == "" {
		return Session{}, errors.New("refresh hash is empty")
	}
	if lifetime < time.Minute {
		return Session{}, errors.New("session lifetime is less than one minute")
	}

	return Session{
		ID:          uuid.New(),
		AccountID:   accountID,
		UserAgent:   userAgent,
		RefreshHash: refreshHash,
		ExpiresAt:   time.Now().UTC().Add(lifetime),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}, nil
}

func (s *Session) Rotate(refreshHash string, lifetime time.Duration) {
	s.RefreshHash = refreshHash
	s.ExpiresAt = time.Now().UTC().Add(lifetime)
}

func (s *Session) IsExpired() bool {
	return time.Now().UTC().After(s.ExpiresAt)
}
