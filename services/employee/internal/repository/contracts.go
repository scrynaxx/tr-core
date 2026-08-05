package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/scrynaxx/tr-core/services/employee/internal/model"
)

// Employee описывает хранение и управление рабочими профилями сотрудников.
type Employee interface {
	// Get возвращает сотрудника по идентификатору, включая архивного.
	Get(ctx context.Context, employeeID uuid.UUID) (model.Employee, error)

	// List возвращает активных сотрудников.
	List(ctx context.Context) ([]model.Employee, error)

	// Create сохраняет нового сотрудника.
	Create(ctx context.Context, employee model.Employee) (model.Employee, error)

	// Save изменяет активного сотрудника.
	Save(ctx context.Context, employee model.Employee) (model.Employee, error)

	// Archive исключает сотрудника из активного списка.
	Archive(ctx context.Context, employeeID uuid.UUID) error
}

// EmployeeCredentials описывает хранение и управление учётными данными сотрудников.
type EmployeeCredentials interface {
	// Get ищет учётные данные сотрудника.
	Get(ctx context.Context, employeeID uuid.UUID) (model.EmployeeCredentials, error)

	// Find ищет учётные данные сотрудника.
	Find(ctx context.Context, employeeID uuid.UUID) (*model.EmployeeCredentials, error)

	// FindByEmail ищет учётные данные сотрудника по email адресу.
	FindByEmail(ctx context.Context, email string) (uuid.UUID, *model.EmployeeCredentials, error)

	// Create создаёт учетные данные сотрудника.
	Create(ctx context.Context, employeeID uuid.UUID, credentials model.EmployeeCredentials) (model.EmployeeCredentials, error)

	// Save обновляет учётные данные сотрудника.
	Save(ctx context.Context, employeeID uuid.UUID, credentials model.EmployeeCredentials) (model.EmployeeCredentials, error)

	// Delete удаляет учётные данные сотрудника.
	Delete(ctx context.Context, employeeID uuid.UUID) error
}

// EmployeePassport описывает хранение и управление паспортными данными сотрудников.
type EmployeePassport interface {
	// Get ищет паспортные данные сотрудника.
	Get(ctx context.Context, employeeID uuid.UUID) (model.EmployeePassport, error)

	// Find ищет паспортные данные сотрудника.
	Find(ctx context.Context, employeeID uuid.UUID) (*model.EmployeePassport, error)

	// Create создаёт паспортные данные сотрудника.
	Create(ctx context.Context, employeeID uuid.UUID, passport model.EmployeePassport) (model.EmployeePassport, error)

	// Save обновляет паспортные данные сотрудника.
	Save(ctx context.Context, employeeID uuid.UUID, passport model.EmployeePassport) (model.EmployeePassport, error)

	// Delete удаляет паспортные данные сотрудника.
	Delete(ctx context.Context, employeeID uuid.UUID) error
}
