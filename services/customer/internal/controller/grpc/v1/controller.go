package v1

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	pbv1 "github.com/scrynaxx/tr-core/contracts/generated/services/v1/customer"
	customerUseCase "github.com/scrynaxx/tr-core/services/customer/internal/model"
	"github.com/scrynaxx/tr-core/services/customer/internal/usecase"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Controller struct {
	customerCase usecase.Customer
}

func NewController(customerCase usecase.Customer) pbv1.CustomerServiceServer {
	return &Controller{customerCase: customerCase}
}

func (c *Controller) CreateCustomer(ctx context.Context, req *pbv1.CreateCustomerRequest) (*pbv1.Customer, error) {
	params := customerUseCase.CreateCustomerParams{
		Name:  req.Name,
		Phone: req.Phone,
		Email: req.Email,
	}

	customer, err := c.customerCase.Create(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("create customer: %w", err)
	}

	return toCustomer(customer), nil
}

func (c *Controller) GetCustomer(ctx context.Context, req *pbv1.GetCustomerRequest) (*pbv1.Customer, error) {
	customer, err := c.customerCase.Get(ctx, uuid.MustParse(req.CustomerId))
	if err != nil {
		return nil, fmt.Errorf("get customer: %w", err)
	}

	return toCustomer(customer), nil
}

func (c *Controller) ListCustomers(ctx context.Context, _ *emptypb.Empty) (*pbv1.ListCustomersResponse, error) {
	customers, err := c.customerCase.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list customers: %w", err)
	}

	return toListCustomersResponse(customers), nil
}

func (c *Controller) UpdateCustomer(ctx context.Context, req *pbv1.UpdateCustomerRequest) (*pbv1.Customer, error) {
	params := customerUseCase.UpdateCustomerParams{
		Name:  req.Name,
		Phone: req.Phone,
		Email: req.Email,
	}

	customer, err := c.customerCase.Update(ctx, uuid.MustParse(req.CustomerId), params)
	if err != nil {
		return nil, fmt.Errorf("update customer: %w", err)
	}

	return toCustomer(customer), nil
}

func (c *Controller) ArchiveCustomer(ctx context.Context, req *pbv1.ArchiveCustomerRequest) (*emptypb.Empty, error) {
	if err := c.customerCase.Archive(ctx, uuid.MustParse(req.CustomerId)); err != nil {
		return nil, fmt.Errorf("archive customer: %w", err)
	}

	return new(emptypb.Empty), nil
}
