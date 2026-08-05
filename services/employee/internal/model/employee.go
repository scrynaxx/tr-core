package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Employee описывает рабочий профиль учётной записи сотрудника.
type Employee struct {
	ID         uuid.UUID    `json:"id"`
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

type CreateEmployeeParams struct {
	Type       EmployeeType
	FirstName  string
	LastName   string
	Patronymic string
	Phone      string
	BirthDate  time.Time

	Credentials *CreateCredentialsParams
	Passport    *CreatePassportParams
}

type UpdateEmployeeParams struct {
	Type       EmployeeType
	FirstName  string
	LastName   string
	Patronymic string
	Phone      string
	BirthDate  time.Time
}

func NewEmployee(p CreateEmployeeParams) (Employee, error) {
	if p.FirstName == "" {
		return Employee{}, errors.New("first name is required")
	}
	if p.LastName == "" {
		return Employee{}, errors.New("last name is required")
	}
	if p.Phone == "" {
		return Employee{}, errors.New("phone is required")
	}
	if p.BirthDate.IsZero() {
		return Employee{}, errors.New("birth date is required")
	}
	if p.Type == EmployeeTypeOwner {
		return Employee{}, errors.New("cannot set employee type to owner")
	}

	return Employee{
		ID:         uuid.New(),
		Type:       p.Type,
		FirstName:  p.FirstName,
		LastName:   p.LastName,
		Patronymic: p.Patronymic,
		Phone:      p.Phone,
		BirthDate:  p.BirthDate,
	}, nil
}

func (e *Employee) Update(p UpdateEmployeeParams) error {
	if p.FirstName == "" {
		return errors.New("first name is required")
	}
	if p.LastName == "" {
		return errors.New("last name is required")
	}
	if p.Phone == "" {
		return errors.New("phone is required")
	}
	if p.BirthDate.IsZero() {
		return errors.New("birth date is required")
	}
	if p.Type == EmployeeTypeOwner {
		return errors.New("cannot set employee type to owner")
	}

	e.Type = p.Type
	e.FirstName = p.FirstName
	e.LastName = p.LastName
	e.Patronymic = p.Patronymic
	e.Phone = p.Phone
	e.BirthDate = p.BirthDate

	return nil
}
