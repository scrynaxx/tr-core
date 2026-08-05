package model

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

const MinPasswordLength = 8

// EmployeeCredentials описывает данные авторизации сотрудника.
type EmployeeCredentials struct {
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateCredentialsParams struct {
	Email    string
	Password string
}

func NewEmployeeCredentials(email, password string) (EmployeeCredentials, error) {
	if utf8.RuneCountInString(strings.TrimSpace(password)) < MinPasswordLength {
		return EmployeeCredentials{}, errors.New("password is required")
	}

	b, err := bcrypt.GenerateFromPassword([]byte(strings.TrimSpace(password)), bcrypt.DefaultCost)
	if err != nil {
		return EmployeeCredentials{}, fmt.Errorf("generate bcrypt hash: %w", err)
	}

	return EmployeeCredentials{
		Email:        email,
		PasswordHash: string(b),
	}, nil
}

func (c *EmployeeCredentials) Update(email string) error {
	if email == "" {
		return errors.New("email is required")
	}

	c.Email = email

	return nil
}
