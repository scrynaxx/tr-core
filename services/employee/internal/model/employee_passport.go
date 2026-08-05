package model

import (
	"errors"
	"time"
)

// EmployeePassport описывает паспортные данные сотрудника.
type EmployeePassport struct {
	Series         string    `json:"series"`
	Number         string    `json:"number"`
	IssuedBy       string    `json:"issued_by"`
	IssuedAt       time.Time `json:"issued_at"`
	DepartmentCode string    `json:"department_code"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CreatePassportParams struct {
	Series         string
	Number         string
	IssuedBy       string
	IssuedAt       time.Time
	DepartmentCode string
}

type UpdatePassportParams struct {
	Series         string
	Number         string
	IssuedBy       string
	IssuedAt       time.Time
	DepartmentCode string
}

func NewEmployeePassport(p CreatePassportParams) (EmployeePassport, error) {
	if p.Series == "" {
		return EmployeePassport{}, errors.New("series is required")
	}
	if p.Number == "" {
		return EmployeePassport{}, errors.New("number is required")
	}
	if p.IssuedBy == "" {
		return EmployeePassport{}, errors.New("issued by is required")
	}
	if p.IssuedAt.IsZero() {
		return EmployeePassport{}, errors.New("issued at time is required")
	}
	if p.DepartmentCode == "" {
		return EmployeePassport{}, errors.New("department is required")
	}

	return EmployeePassport{
		Series:         p.Series,
		Number:         p.Number,
		IssuedBy:       p.IssuedBy,
		IssuedAt:       p.IssuedAt,
		DepartmentCode: p.DepartmentCode,
	}, nil
}

func (p *EmployeePassport) Update(pa UpdatePassportParams) error {
	if pa.Series == "" {
		return errors.New("series is required")
	}
	if pa.Number == "" {
		return errors.New("number is required")
	}
	if pa.IssuedBy == "" {
		return errors.New("issued by is required")
	}
	if pa.IssuedAt.IsZero() {
		return errors.New("issued at time is required")
	}
	if pa.DepartmentCode == "" {
		return errors.New("department is required")
	}

	p.Series = pa.Series
	p.Number = pa.Number
	p.IssuedBy = pa.IssuedBy
	p.IssuedAt = pa.IssuedAt
	p.DepartmentCode = pa.DepartmentCode

	return nil
}
