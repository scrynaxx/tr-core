package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/scrynaxx/tr-core/services/employee/internal/model"
)

// Employee управляет рабочими профилями сотрудников.
type Employee interface {
	// Create создаёт сотрудника для существующей учётной записи.
	Create(ctx context.Context, p model.CreateEmployeeParams) (model.Employee, error)

	// Get возвращает сотрудника по идентификатору, включая архивного.
	Get(ctx context.Context, employeeID uuid.UUID) (model.Employee, error)

	// List возвращает активных сотрудников.
	List(ctx context.Context) ([]model.Employee, error)

	// Update изменяет активного сотрудника.
	Update(ctx context.Context, employeeID uuid.UUID, p model.UpdateEmployeeParams) (model.Employee, error)

	// Archive исключает сотрудника из дальнейшего назначения. Возвращает ошибку в случае попытки архивировать себя или сотрудника с типом owner.
	Archive(ctx context.Context, callerID, employeeID uuid.UUID) error

	// FindCredentials ищет учётные данные сотрудника.
	FindCredentials(ctx context.Context, employeeID uuid.UUID) (*model.EmployeeCredentials, error)

	// FindCredentialsByEmail ищет учётные данные сотрудника по email адресу.
	FindCredentialsByEmail(ctx context.Context, email string) (uuid.UUID, *model.EmployeeCredentials, error)

	// CreateCredentials создаёт учетные данные сотрудника.
	CreateCredentials(ctx context.Context, employeeID uuid.UUID, p model.CreateCredentialsParams) (model.EmployeeCredentials, error)

	// UpdateCredentials обновляет учётные данные сотрудника
	UpdateCredentials(ctx context.Context, employeeID uuid.UUID, email string) (model.EmployeeCredentials, error)

	// DeleteCredentials удаляет учётные данные сотрудника.
	DeleteCredentials(ctx context.Context, employeeID uuid.UUID) error

	// FindPassport ищет паспортные данные сотрудника
	FindPassport(ctx context.Context, employeeID uuid.UUID) (*model.EmployeePassport, error)

	// CreatePassport создаёт паспортные данные сотрудника.
	CreatePassport(ctx context.Context, employeeID uuid.UUID, p model.CreatePassportParams) (model.EmployeePassport, error)

	// UpdatePassport обновляет паспортные данные сотрудника.
	UpdatePassport(ctx context.Context, employeeID uuid.UUID, p model.UpdatePassportParams) (model.EmployeePassport, error)

	// DeletePassport удаляет паспортные данные сотрудника.
	DeletePassport(ctx context.Context, employeeID uuid.UUID) error
}
