package model

import (
	"errors"
	"time"
)

type Passport struct {
	Series         string    `json:"series"`
	Number         string    `json:"number"`
	IssuedBy       string    `json:"issued_by"`
	IssuedAt       time.Time `json:"issued_at"`
	DepartmentCode string    `json:"department_code"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type PassportInput struct {
	Series         string
	Number         string
	IssuedBy       string
	IssuedAt       time.Time
	DepartmentCode string
}

func NewPassport(input PassportInput) (Passport, error) {
	if err := input.Validate(); err != nil {
		return Passport{}, err
	}

	return Passport{
		Series:         input.Series,
		Number:         input.Number,
		IssuedBy:       input.IssuedBy,
		IssuedAt:       input.IssuedAt,
		DepartmentCode: input.DepartmentCode,
	}, nil
}

func (p *Passport) Update(input PassportInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	p.Series = input.Series
	p.Number = input.Number
	p.IssuedBy = input.IssuedBy
	p.IssuedAt = input.IssuedAt
	p.DepartmentCode = input.DepartmentCode

	return nil
}

func (i PassportInput) Validate() error {
	if i.Series == "" {
		return errors.New("passport series is required")
	}
	if i.Number == "" {
		return errors.New("passport number is required")
	}
	if i.IssuedBy == "" {
		return errors.New("passport issuer is required")
	}
	if i.IssuedAt.IsZero() {
		return errors.New("passport issue date is required")
	}
	if i.DepartmentCode == "" {
		return errors.New("passport department code is required")
	}

	return nil
}
