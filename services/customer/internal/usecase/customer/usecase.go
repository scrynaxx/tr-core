package customer

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/scrynaxx/tr-core/services/customer/internal/model"
	"github.com/scrynaxx/tr-core/services/customer/internal/repository"
	"github.com/scrynaxx/tr-core/services/customer/internal/usecase"
)

type UseCase struct {
	customerRepo repository.Customer
}

func NewUseCase(customerRepo repository.Customer) usecase.Customer {
	return withTracing(&UseCase{
		customerRepo: customerRepo,
	})
}

func (u *UseCase) Create(ctx context.Context, params model.CreateCustomerParams) (model.Customer, error) {
	customer, err := model.NewCustomer(params)
	if err != nil {
		return model.Customer{}, fmt.Errorf("new customer: %w", err)
	}

	customer, err = u.customerRepo.Create(ctx, customer)
	if err != nil {
		return model.Customer{}, fmt.Errorf("insert customer: %w", err)
	}

	return customer, nil
}

func (u *UseCase) Get(ctx context.Context, customerID uuid.UUID) (model.Customer, error) {
	customer, err := u.customerRepo.Get(ctx, customerID)
	if err != nil {
		return model.Customer{}, fmt.Errorf("get customer: %w", err)
	}

	return customer, nil
}

func (u *UseCase) List(ctx context.Context) ([]model.Customer, error) {
	customers, err := u.customerRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list customers: %w", err)
	}

	return customers, nil
}

func (u *UseCase) Update(ctx context.Context, customerID uuid.UUID, params model.UpdateCustomerParams) (model.Customer, error) {
	customer, err := u.customerRepo.Get(ctx, customerID)
	if err != nil {
		return model.Customer{}, fmt.Errorf("get customer: %w", err)
	}

	if err = customer.Update(params); err != nil {
		return model.Customer{}, fmt.Errorf("update customer: %w", err)
	}

	customer, err = u.customerRepo.Save(ctx, customer)
	if err != nil {
		return model.Customer{}, fmt.Errorf("save customer: %w", err)
	}

	return customer, nil
}

func (u *UseCase) Archive(ctx context.Context, customerID uuid.UUID) error {
	if err := u.customerRepo.Archive(ctx, customerID); err != nil {
		return fmt.Errorf("archive customer: %w", err)
	}

	return nil
}
