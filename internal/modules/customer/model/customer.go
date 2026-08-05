package model

import (
	"errors"
	"time"

	"uuid"

	"github.com/samber/lo"
)

// Customer заказчик перевозки.
type Customer struct {
	ID         uuid.UUID  `json:"id"`
	AccountID  *uuid.UUID `json:"account_id"`
	FirstName  string     `json:"first_name"`
	LastName   string     `json:"last_name"`
	Patronymic *string    `json:"patronymic"`
	Phone      string     `json:"phone"`
	Email      string     `json:"email"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	ArchivedAt *time.Time `json:"archived_at"`
}

type CustomerInput struct {
	FirstName  string
	LastName   string
	Patronymic *string
	Phone      string
	Email      string
}

func NewCustomer(input CustomerInput) (Customer, error) {
	if input.FirstName == "" {
		return Customer{}, errors.New("first name is required")
	}
	if input.LastName == "" {
		return Customer{}, errors.New("last name is required")
	}
	if input.Email == "" {
		return Customer{}, errors.New("email is required")
	}
	if input.Phone == "" {
		return Customer{}, errors.New("phone is required")
	}

	return Customer{
		ID:         uuid.New(),
		FirstName:  input.FirstName,
		LastName:   input.LastName,
		Patronymic: input.Patronymic,
		Phone:      input.Phone,
		Email:      input.Email,
	}, nil
}

func (c *Customer) Update(input CustomerInput) error {
	if input.FirstName == "" {
		return errors.New("first name is required")
	}
	if input.LastName == "" {
		return errors.New("last name is required")
	}
	if input.Email == "" {
		return errors.New("email is required")
	}
	if input.Phone == "" {
		return errors.New("phone is required")
	}

	c.FirstName = input.FirstName
	c.LastName = input.LastName
	c.Patronymic = input.Patronymic
	c.Phone = input.Phone
	c.Email = input.Email

	return nil
}

func (c *Customer) Archive() {
	if c.ArchivedAt != nil {
		return
	}

	c.ArchivedAt = lo.ToPtr(time.Now().UTC())
}

func (c *Customer) Restore() {
	c.ArchivedAt = nil
}
