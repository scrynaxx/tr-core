package contract

import (
	"time"

	"uuid"
)

type CreateCustomerRequest struct {
	FirstName  string  `json:"first_name" validate:"required"`
	LastName   string  `json:"last_name" validate:"required"`
	Patronymic *string `json:"patronymic"`
	Phone      string  `json:"phone" validate:"required"`
	Email      string  `json:"email" validate:"required"`
}

type GetCustomerRequest struct {
	CustomerID uuid.UUID `param:"customer_id" validate:"required,uuid"`
}

type UpdateCustomerRequest struct {
	CustomerID uuid.UUID `param:"customer_id" validate:"required,uuid"`
	FirstName  string    `json:"first_name" validate:"required"`
	LastName   string    `json:"last_name" validate:"required"`
	Patronymic *string   `json:"patronymic"`
	Phone      string    `json:"phone" validate:"required"`
	Email      string    `json:"email" validate:"required"`
}

type ArchiveCustomerRequest struct {
	CustomerID uuid.UUID `param:"customer_id" validate:"required,uuid"`
}

type RestoreCustomerRequest struct {
	CustomerID uuid.UUID `param:"customer_id" validate:"required,uuid"`
}

type ListCustomersResponse struct {
	Customers []Customer `json:"customers"`
}

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
