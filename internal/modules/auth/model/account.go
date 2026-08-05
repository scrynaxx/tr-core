package model

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
	"uuid"

	"github.com/scrynaxx/tr-core/internal/security"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	ID           uuid.UUID      `json:"id"`
	Actor        security.Actor `json:"actor"`
	Email        string         `json:"email"`
	PasswordHash string         `json:"password_hash"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

func NewAccount(actor security.Actor, email string, password string) (Account, error) {
	if email == "" {
		return Account{}, errors.New("email is required")
	}
	if password == "" {
		return Account{}, errors.New("password is required")
	}

	account := Account{
		ID:           uuid.New(),
		Actor:        actor,
		Email:        email,
		PasswordHash: "",
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	return account, account.SetPassword(password)
}

func (a *Account) SetPassword(password string) error {
	const MinPasswordLength = 8

	password = strings.TrimSpace(password)
	if utf8.RuneCountInString(password) < MinPasswordLength {
		return fmt.Errorf("password is too short, minimum length is %d characters", MinPasswordLength)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	a.PasswordHash = string(hash)

	return nil
}

func (a *Account) PasswordEqual(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(a.PasswordHash), []byte(password)) == nil
}
