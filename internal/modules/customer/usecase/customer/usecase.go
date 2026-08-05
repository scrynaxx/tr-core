package customer

import (
	"context"
	"uuid"

	"github.com/scrynaxx/tr-core/internal/modules/customer/model"
	"github.com/scrynaxx/tr-core/internal/modules/customer/repository"
	"github.com/scrynaxx/tr-core/internal/modules/customer/usecase"
)

type UseCase struct {
	customerRepo repository.Customer
}

func NewUseCase(customerRepo repository.Customer) usecase.Customer {
	return &UseCase{
		customerRepo: customerRepo,
	}
}

func (uc *UseCase) Create(ctx context.Context, input model.CustomerInput) error {
	customer, err := model.NewCustomer(input)
	if err != nil {
		return err
	}

	return uc.customerRepo.Create(ctx, customer)
}

func (uc *UseCase) Get(ctx context.Context, customerID uuid.UUID) (model.Customer, error) {
	return uc.customerRepo.Get(ctx, customerID)
}

func (uc *UseCase) List(ctx context.Context) ([]model.Customer, error) {
	return uc.customerRepo.List(ctx)
}

func (uc *UseCase) Update(ctx context.Context, customerID uuid.UUID, input model.CustomerInput) error {
	customer, err := uc.customerRepo.Get(ctx, customerID)
	if err != nil {
		return err
	}

	if err = customer.Update(input); err != nil {
		return err
	}

	return uc.customerRepo.Update(ctx, customer)
}

func (uc *UseCase) Archive(ctx context.Context, customerID uuid.UUID) error {
	customer, err := uc.customerRepo.Get(ctx, customerID)
	if err != nil {
		return err
	}

	customer.Archive()

	return uc.customerRepo.Update(ctx, customer)
}

func (uc *UseCase) Restore(ctx context.Context, customerID uuid.UUID) error {
	customer, err := uc.customerRepo.Get(ctx, customerID)
	if err != nil {
		return err
	}

	customer.Restore()

	return uc.customerRepo.Update(ctx, customer)
}
