package employee

import (
	"context"
	"uuid"

	"github.com/scrynaxx/tr-core/internal/event"
	"github.com/scrynaxx/tr-core/internal/modules/employee/model"
	"github.com/scrynaxx/tr-core/internal/modules/employee/repository"
	"github.com/scrynaxx/tr-core/internal/modules/employee/usecase"
	"github.com/scrynaxx/tr-core/pkg/database"
	"github.com/scrynaxx/tr-core/pkg/messaging"
)

type UseCase struct {
	employeeRepo repository.Employee
	passportRepo repository.Passport
	outboxRepo   messaging.OutboxRepository
	transactor   database.Transactor
}

func New(
	employeeRepo repository.Employee,
	passportRepo repository.Passport,
	outboxRepo messaging.OutboxRepository,
	transactor database.Transactor) usecase.Employee {
	return &UseCase{
		employeeRepo: employeeRepo,
		passportRepo: passportRepo,
		outboxRepo:   outboxRepo,
		transactor:   transactor,
	}
}

func (uc *UseCase) Create(ctx context.Context, input model.EmployeeInput, passportInput *model.PassportInput) error {
	employee, err := model.NewEmployee(input)
	if err != nil {
		return err
	}

	if passportInput == nil {
		return uc.employeeRepo.Create(ctx, employee)
	}

	passport, err := model.NewPassport(*passportInput)
	if err != nil {
		return err
	}

	return database.Tx(ctx, uc.transactor, func(ctx context.Context) error {
		if err = uc.employeeRepo.Create(ctx, employee); err != nil {
			return err
		}

		return uc.passportRepo.Create(ctx, employee.ID, passport)
	})
}

func (uc *UseCase) Get(ctx context.Context, employeeID uuid.UUID) (model.Employee, error) {
	return uc.employeeRepo.Get(ctx, employeeID)
}

func (uc *UseCase) GetByAccount(ctx context.Context, accountID uuid.UUID) (model.Employee, error) {
	return uc.employeeRepo.GetByAccount(ctx, accountID)
}

func (uc *UseCase) List(ctx context.Context) ([]model.Employee, error) {
	return uc.employeeRepo.List(ctx)
}

func (uc *UseCase) Update(ctx context.Context, employeeID uuid.UUID, input model.EmployeeInput) error {
	employee, err := uc.employeeRepo.Get(ctx, employeeID)
	if err != nil {
		return err
	}

	if err = employee.Update(input); err != nil {
		return err
	}

	return uc.employeeRepo.Update(ctx, employee)
}

func (uc *UseCase) Archive(ctx context.Context, accountID, employeeID uuid.UUID) error {
	employee, err := uc.employeeRepo.Get(ctx, employeeID)
	if err != nil {
		return err
	}

	if err = employee.Archive(accountID); err != nil {
		return err
	}

	return database.Tx(ctx, uc.transactor, func(ctx context.Context) error {
		if err = uc.employeeRepo.Update(ctx, employee); err != nil {
			return err
		}

		msg, err := messaging.NewMessage(event.EmployeeArchivedV1, event.EmployeeArchivedV1Payload{EmployeeID: employeeID})
		if err != nil {
			return err
		}

		return uc.outboxRepo.StoreEvent(ctx, msg)
	})
}

func (uc *UseCase) Restore(ctx context.Context, employeeID uuid.UUID) error {
	employee, err := uc.employeeRepo.Get(ctx, employeeID)
	if err != nil {
		return err
	}

	employee.Restore()

	return uc.employeeRepo.Update(ctx, employee)
}
