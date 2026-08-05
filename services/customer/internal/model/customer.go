package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Customer заказчик перевозки.
type Customer struct {
	ID         uuid.UUID  `json:"id"`
	Name       string     `json:"name"`
	Phone      string     `json:"phone"`
	Email      string     `json:"email"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	ArchivedAt *time.Time `json:"archived_at"`
}

type CreateCustomerParams struct {
	Name  string
	Phone string
	Email string
}

type UpdateCustomerParams struct {
	Name  string
	Phone string
	Email string
}

func NewCustomer(params CreateCustomerParams) (Customer, error) {
	if params.Name == "" {
		return Customer{}, errors.New("name is required")
	}
	if params.Email == "" {
		return Customer{}, errors.New("email is required")
	}
	if params.Phone == "" {
		return Customer{}, errors.New("phone is required")
	}

	return Customer{
		ID:    uuid.New(),
		Name:  params.Name,
		Phone: params.Phone,
		Email: params.Email,
	}, nil
}

func (c *Customer) Update(params UpdateCustomerParams) error {
	if params.Name == "" {
		return errors.New("name is required")
	}
	if params.Email == "" {
		return errors.New("email is required")
	}
	if params.Phone == "" {
		return errors.New("phone is required")
	}

	c.Name = params.Name
	c.Phone = params.Phone
	c.Email = params.Email

	return nil
}
