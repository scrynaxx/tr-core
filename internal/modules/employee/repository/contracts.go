package repository

import (
	"context"
	"uuid"

	"github.com/scrynaxx/tr-core/internal/modules/employee/model"
)

type Employee interface {
	Get(ctx context.Context, employeeID uuid.UUID) (model.Employee, error)
	GetByAccount(ctx context.Context, accountID uuid.UUID) (model.Employee, error)
	List(ctx context.Context) ([]model.Employee, error)
	Create(ctx context.Context, employee model.Employee) error
	Update(ctx context.Context, employee model.Employee) error
}

type Passport interface {
	Get(ctx context.Context, employeeID uuid.UUID) (model.Passport, error)
	Find(ctx context.Context, employeeID uuid.UUID) (*model.Passport, error)
	Create(ctx context.Context, employeeID uuid.UUID, passport model.Passport) error
	Update(ctx context.Context, employeeID uuid.UUID, passport model.Passport) error
	Delete(ctx context.Context, employeeID uuid.UUID) error
}
