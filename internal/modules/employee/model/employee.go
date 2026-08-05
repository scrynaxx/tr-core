package model

import (
	"errors"
	"time"

	"uuid"

	"github.com/samber/lo"
)

type Employee struct {
	ID         uuid.UUID    `json:"id"`
	AccountID  *uuid.UUID   `json:"account_id"`
	Type       EmployeeType `json:"type"`
	FirstName  string       `json:"first_name"`
	LastName   string       `json:"last_name"`
	Patronymic string       `json:"patronymic"`
	Phone      string       `json:"phone"`
	BirthDate  time.Time    `json:"birth_date"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
	ArchivedAt *time.Time   `json:"archived_at"`
}

type EmployeeInput struct {
	Type       EmployeeType
	FirstName  string
	LastName   string
	Patronymic string
	Phone      string
	BirthDate  time.Time
}

func NewEmployee(input EmployeeInput) (Employee, error) {
	if err := input.Validate(); err != nil {
		return Employee{}, err
	}

	return Employee{
		ID:         uuid.New(),
		Type:       input.Type,
		FirstName:  input.FirstName,
		LastName:   input.LastName,
		Patronymic: input.Patronymic,
		Phone:      input.Phone,
		BirthDate:  input.BirthDate,
	}, nil
}

func (e *Employee) Update(input EmployeeInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	e.Type = input.Type
	e.FirstName = input.FirstName
	e.LastName = input.LastName
	e.Patronymic = input.Patronymic
	e.Phone = input.Phone
	e.BirthDate = input.BirthDate

	return nil
}

func (e *Employee) Archive(accountID uuid.UUID) error {
	if e.AccountID != nil && *e.AccountID == accountID {
		return errors.New("cannot archive yourself")
	}
	if e.Type == EmployeeTypeOwner {
		return errors.New("cannot archive owner")
	}

	e.ArchivedAt = lo.ToPtr(time.Now().UTC())

	return nil
}

func (e *Employee) Restore() {
	e.ArchivedAt = nil
}

func (i EmployeeInput) Validate() error {
	if i.FirstName == "" {
		return errors.New("first name is required")
	}
	if i.LastName == "" {
		return errors.New("last name is required")
	}
	if i.Phone == "" {
		return errors.New("phone is required")
	}
	if i.BirthDate.IsZero() {
		return errors.New("birth date is required")
	}

	return nil
}
