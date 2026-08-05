package inner

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	pbinner "github.com/scrynaxx/tr-core/contracts/generated/services/inner/customer"
	"github.com/scrynaxx/tr-core/services/customer/internal/usecase"
)

type Controller struct {
	customerCase usecase.Customer
}

func NewController(customerCase usecase.Customer) pbinner.CustomerServiceServer {
	return &Controller{customerCase: customerCase}
}

func (c *Controller) GetCustomer(ctx context.Context, req *pbinner.GetCustomerRequest) (*pbinner.GetCustomerResponse, error) {
	customerID, err := uuid.Parse(req.CustomerId)
	if err != nil {
		return nil, fmt.Errorf("parse customer id: %w", err)
	}

	customer, err := c.customerCase.Get(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("get customer: %w", err)
	}

	return toGetCustomerResponse(customer), nil
}
