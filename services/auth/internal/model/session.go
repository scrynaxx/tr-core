package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID          uuid.UUID `json:"id"`
	EmployeeID  uuid.UUID `json:"employee_id"`
	UserAgent   string    `json:"user_agent"`
	RefreshHash string    `json:"refresh_hash"`
	ExpiresAt   time.Time `json:"expires_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewSession(employeeID uuid.UUID, userAgent, refreshHash string, lifetime time.Duration) (Session, error) {
	if employeeID == uuid.Nil {
		return Session{}, errors.New("employee ID is nil")
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

	now := time.Now().UTC()

	return Session{
		ID:          uuid.New(),
		EmployeeID:  employeeID,
		UserAgent:   userAgent,
		RefreshHash: refreshHash,
		ExpiresAt:   now.Add(lifetime),
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (s *Session) UpdateRefresh(hash string, lifetime time.Duration) {
	s.ExpiresAt = time.Now().UTC().Add(lifetime)
	s.RefreshHash = hash
}

func (s *Session) IsExpired() bool {
	return s.ExpiresAt.Before(time.Now().UTC())
}
