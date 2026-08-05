package usecase

import (
	"context"
	"uuid"

	"github.com/scrynaxx/tr-core/internal/modules/employee/model"
)

type Employee interface {
	Create(ctx context.Context, input model.EmployeeInput, passportInput *model.PassportInput) error
	Get(ctx context.Context, employeeID uuid.UUID) (model.Employee, error)
	GetByAccount(ctx context.Context, accountID uuid.UUID) (model.Employee, error)
	List(ctx context.Context) ([]model.Employee, error)
	Update(ctx context.Context, employeeID uuid.UUID, input model.EmployeeInput) error
	Archive(ctx context.Context, accountID, employeeID uuid.UUID) error
	Restore(ctx context.Context, employeeID uuid.UUID) error
}

type Passport interface {
	Find(ctx context.Context, employeeID uuid.UUID) (*model.Passport, error)
	Create(ctx context.Context, employeeID uuid.UUID, input model.PassportInput) error
	Update(ctx context.Context, employeeID uuid.UUID, input model.PassportInput) error
	Delete(ctx context.Context, employeeID uuid.UUID) error
}
