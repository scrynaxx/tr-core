package employee

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/scrynaxx/tr-core/contracts/globalevent"
	"github.com/scrynaxx/tr-core/pkg/database"
	"github.com/scrynaxx/tr-core/pkg/events"
	"github.com/scrynaxx/tr-core/services/employee/internal/model"
	"github.com/scrynaxx/tr-core/services/employee/internal/repository"
	"github.com/scrynaxx/tr-core/services/employee/internal/usecase"
)

type UseCase struct {
	employeeRepo    repository.Employee
	credentialsRepo repository.EmployeeCredentials
	passportRepo    repository.EmployeePassport
	eventRepo       events.OutboxRepository
	transactor      database.Transactor
}

func NewUseCase(
	employeeRepo repository.Employee,
	credentialsRepo repository.EmployeeCredentials,
	passportRepo repository.EmployeePassport,
	eventRepo events.OutboxRepository,
	transactor database.Transactor,
) usecase.Employee {
	return withTracing(&UseCase{
		employeeRepo:    employeeRepo,
		credentialsRepo: credentialsRepo,
		passportRepo:    passportRepo,
		eventRepo:       eventRepo,
		transactor:      transactor,
	})
}

func (u *UseCase) Create(ctx context.Context, p model.CreateEmployeeParams) (model.Employee, error) {
	employee, err := model.NewEmployee(p)
	if err != nil {
		return model.Employee{}, fmt.Errorf("new employee: %w", err)
	}

	return database.WithTxResult(ctx, u.transactor, func(ctx context.Context) (model.Employee, error) {
		employee, err = u.employeeRepo.Create(ctx, employee)
		if err != nil {
			return model.Employee{}, fmt.Errorf("create employee: %w", err)
		}

		if p.Passport != nil {
			passport, err := model.NewEmployeePassport(*p.Passport)
			if err != nil {
				return model.Employee{}, fmt.Errorf("new passport: %w", err)
			}

			if passport, err = u.passportRepo.Create(ctx, employee.ID, passport); err != nil {
				return model.Employee{}, fmt.Errorf("create passport: %w", err)
			}
		}

		if p.Credentials != nil {
			credentials, err := model.NewEmployeeCredentials(p.Credentials.Email, p.Credentials.Password)
			if err != nil {
				return model.Employee{}, fmt.Errorf("new credentials: %w", err)
			}

			if credentials, err = u.credentialsRepo.Create(ctx, employee.ID, credentials); err != nil {
				return model.Employee{}, fmt.Errorf("create credentials: %w", err)
			}
		}

		return employee, nil
	})
}

func (u *UseCase) Get(ctx context.Context, employeeID uuid.UUID) (model.Employee, error) {
	employee, err := u.employeeRepo.Get(ctx, employeeID)
	if err != nil {
		return model.Employee{}, fmt.Errorf("get employee: %w", err)
	}

	return employee, nil
}

func (u *UseCase) List(ctx context.Context) ([]model.Employee, error) {
	employees, err := u.employeeRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list employees: %w", err)
	}

	return employees, nil
}

func (u *UseCase) Update(ctx context.Context, employeeID uuid.UUID, p model.UpdateEmployeeParams) (model.Employee, error) {
	employee, err := u.employeeRepo.Get(ctx, employeeID)
	if err != nil {
		return model.Employee{}, fmt.Errorf("get employee: %w", err)
	}

	if err = employee.Update(p); err != nil {
		return model.Employee{}, fmt.Errorf("update employee: %w", err)
	}

	employee, err = u.employeeRepo.Save(ctx, employee)
	if err != nil {
		return model.Employee{}, fmt.Errorf("save employee: %w", err)
	}

	return employee, nil
}

func (u *UseCase) Archive(ctx context.Context, callerID, employeeID uuid.UUID) error {
	employee, err := u.employeeRepo.Get(ctx, employeeID)
	if err != nil {
		if errors.Is(err, model.ErrEmployeeNotFound) {
			return nil
		}
		return fmt.Errorf("get employee: %w", err)
	}

	if employee.ID == callerID {
		return fmt.Errorf("cannot archive yourself")
	}

	if employee.Type == model.EmployeeTypeOwner {
		return fmt.Errorf("employee type is owner, not allowed")
	}

	return database.WithTx(ctx, u.transactor, func(ctx context.Context) error {
		if err := u.employeeRepo.Archive(ctx, employeeID); err != nil {
			return fmt.Errorf("archive employee: %w", err)
		}

		e, err := events.NewMessage(globalevent.EmployeeArchivedV1, globalevent.EmployeeArchivedDataV1{EmployeeID: employeeID})
		if err != nil {
			return fmt.Errorf("create employee archived event: %w", err)
		}

		if err = u.eventRepo.StoreEvent(ctx, e); err != nil {
			return fmt.Errorf("store employee archived event: %w", err)
		}

		return nil
	})
}

func (u *UseCase) FindCredentials(ctx context.Context, employeeID uuid.UUID) (*model.EmployeeCredentials, error) {
	credentials, err := u.credentialsRepo.Find(ctx, employeeID)
	if err != nil {
		return nil, fmt.Errorf("find credentials: %w", err)
	}

	return credentials, nil
}

func (u *UseCase) FindCredentialsByEmail(ctx context.Context, email string) (uuid.UUID, *model.EmployeeCredentials, error) {
	employeeID, credentials, err := u.credentialsRepo.FindByEmail(ctx, email)
	if err != nil {
		return uuid.Nil, nil, fmt.Errorf("find credentials by email: %w", err)
	}

	return employeeID, credentials, nil
}

func (u *UseCase) CreateCredentials(ctx context.Context, employeeID uuid.UUID, p model.CreateCredentialsParams) (model.EmployeeCredentials, error) {
	credentials, err := model.NewEmployeeCredentials(p.Email, p.Password)
	if err != nil {
		return model.EmployeeCredentials{}, fmt.Errorf("new credentials: %w", err)
	}

	credentials, err = u.credentialsRepo.Create(ctx, employeeID, credentials)
	if err != nil {
		return model.EmployeeCredentials{}, fmt.Errorf("create credentials: %w", err)
	}

	return credentials, nil
}

func (u *UseCase) UpdateCredentials(ctx context.Context, employeeID uuid.UUID, email string) (model.EmployeeCredentials, error) {
	credentials, err := u.credentialsRepo.Get(ctx, employeeID)
	if err != nil {
		return model.EmployeeCredentials{}, fmt.Errorf("get credentials: %w", err)
	}

	if err = credentials.Update(email); err != nil {
		return model.EmployeeCredentials{}, fmt.Errorf("update credentials: %w", err)
	}

	return database.WithTxResult(ctx, u.transactor, func(ctx context.Context) (model.EmployeeCredentials, error) {
		credentials, err = u.credentialsRepo.Save(ctx, employeeID, credentials)
		if err != nil {
			return model.EmployeeCredentials{}, fmt.Errorf("save credentials: %w", err)
		}

		e, err := events.NewMessage(globalevent.EmployeeCredentialsChangedV1, globalevent.EmployeeCredentialsChangedDataV1{EmployeeID: employeeID})
		if err != nil {
			return model.EmployeeCredentials{}, fmt.Errorf("create credentials changed event: %w", err)
		}

		if err = u.eventRepo.StoreEvent(ctx, e); err != nil {
			return model.EmployeeCredentials{}, fmt.Errorf("store credentials changed event: %w", err)
		}

		return credentials, nil
	})
}

func (u *UseCase) DeleteCredentials(ctx context.Context, employeeID uuid.UUID) error {
	return database.WithTx(ctx, u.transactor, func(ctx context.Context) error {
		if err := u.credentialsRepo.Delete(ctx, employeeID); err != nil {
			return fmt.Errorf("delete employee credentials: %w", err)
		}

		e, err := events.NewMessage(globalevent.EmployeeCredentialsChangedV1, globalevent.EmployeeCredentialsChangedDataV1{EmployeeID: employeeID})
		if err != nil {
			return fmt.Errorf("create employee credentials changed event: %w", err)
		}

		if err = u.eventRepo.StoreEvent(ctx, e); err != nil {
			return fmt.Errorf("store employee credentials changed event: %w", err)
		}

		return nil
	})
}

func (u *UseCase) FindPassport(ctx context.Context, employeeID uuid.UUID) (*model.EmployeePassport, error) {
	passport, err := u.passportRepo.Find(ctx, employeeID)
	if err != nil {
		return nil, fmt.Errorf("find passport: %w", err)
	}

	return passport, nil
}

func (u *UseCase) CreatePassport(ctx context.Context, employeeID uuid.UUID, p model.CreatePassportParams) (model.EmployeePassport, error) {
	passport, err := model.NewEmployeePassport(p)
	if err != nil {
		return model.EmployeePassport{}, fmt.Errorf("mew passport: %w", err)
	}

	passport, err = u.passportRepo.Create(ctx, employeeID, passport)
	if err != nil {
		return model.EmployeePassport{}, fmt.Errorf("create passport: %w", err)
	}

	return passport, nil
}

func (u *UseCase) UpdatePassport(ctx context.Context, employeeID uuid.UUID, p model.UpdatePassportParams) (model.EmployeePassport, error) {
	passport, err := u.passportRepo.Get(ctx, employeeID)
	if err != nil {
		return model.EmployeePassport{}, fmt.Errorf("get passport: %w", err)
	}

	if err = passport.Update(p); err != nil {
		return model.EmployeePassport{}, fmt.Errorf("update passport: %w", err)
	}

	passport, err = u.passportRepo.Save(ctx, employeeID, passport)
	if err != nil {
		return model.EmployeePassport{}, fmt.Errorf("save passport: %w", err)
	}

	return passport, nil
}

func (u *UseCase) DeletePassport(ctx context.Context, employeeID uuid.UUID) error {
	if err := u.passportRepo.Delete(ctx, employeeID); err != nil {
		return fmt.Errorf("delete employee passport: %w", err)
	}

	return nil
}
