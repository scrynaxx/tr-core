package contract

import (
	"time"

	"github.com/scrynaxx/tr-core/internal/modules/employee/model"

	"uuid"
)

type CreateEmployeeRequest struct {
	EmployeeRequestBody
	Passport *PassportRequestBody `json:"passport"`
}

type GetEmployeeRequest struct {
	EmployeeID uuid.UUID `param:"employee_id" validate:"required,uuid"`
}

type UpdateEmployeeRequest struct {
	EmployeeRequestBody
	EmployeeID uuid.UUID `param:"employee_id" validate:"required,uuid"`
}

type ArchiveEmployeeRequest struct {
	EmployeeID uuid.UUID `param:"employee_id" validate:"required,uuid"`
}

type RestoreEmployeeRequest struct {
	EmployeeID uuid.UUID `param:"employee_id" validate:"required,uuid"`
}

type CreatePassportRequest struct {
	PassportRequestBody
	EmployeeID uuid.UUID `param:"employee_id" validate:"required,uuid"`
}

type GetPassportRequest struct {
	EmployeeID uuid.UUID `param:"employee_id" validate:"required,uuid"`
}

type UpdatePassportRequest struct {
	PassportRequestBody
	EmployeeID uuid.UUID `json:"employee_id" validate:"required,uuid"`
}

type DeletePassportRequest struct {
	EmployeeID uuid.UUID `param:"employee_id" validate:"required,uuid"`
}

type ListEmployeesResponse struct {
	Employees []Employee `json:"employees"`
}

type EmployeeRequestBody struct {
	Type       model.EmployeeType `json:"type" validate:"required,oneof=manager foreman loader assembler"`
	FirstName  string             `json:"first_name" validate:"required"`
	LastName   string             `json:"last_name" validate:"required"`
	Patronymic string             `json:"patronymic"`
	Phone      string             `json:"phone" validate:"required"`
	BirthDate  time.Time          `json:"birth_date" validate:"required"`
}

type PassportRequestBody struct {
	Series         string    `json:"series" validate:"required"`
	Number         string    `json:"number" validate:"required"`
	IssuedBy       string    `json:"issued_by" validate:"required"`
	IssuedAt       time.Time `json:"issued_at" validate:"required"`
	DepartmentCode string    `json:"department_code" validate:"required"`
}

type Employee struct {
	ID         uuid.UUID          `json:"id"`
	AccountID  *uuid.UUID         `json:"account_id"`
	Type       model.EmployeeType `json:"type"`
	FirstName  string             `json:"first_name"`
	LastName   string             `json:"last_name"`
	Patronymic string             `json:"patronymic"`
	Phone      string             `json:"phone"`
	BirthDate  time.Time          `json:"birth_date"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
	ArchivedAt *time.Time         `json:"archived_at"`
}

type Passport struct {
	Series         string    `json:"series"`
	Number         string    `json:"number"`
	IssuedBy       string    `json:"issued_by"`
	IssuedAt       time.Time `json:"issued_at"`
	DepartmentCode string    `json:"department_code"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Profile struct {
	ID         uuid.UUID          `json:"id"`
	AccountID  uuid.UUID          `json:"account_id"`
	Type       model.EmployeeType `json:"type"`
	FirstName  string             `json:"first_name"`
	LastName   string             `json:"last_name"`
	Patronymic string             `json:"patronymic"`
	Phone      string             `json:"phone"`
	BirthDate  time.Time          `json:"birth_date"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`

	Email string `json:"email"`
}
