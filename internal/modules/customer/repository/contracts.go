package repository

import (
	"context"
	"uuid"

	"github.com/scrynaxx/tr-core/internal/modules/customer/model"
)

type Customer interface {
	Get(ctx context.Context, customerID uuid.UUID) (model.Customer, error)
	List(ctx context.Context) ([]model.Customer, error)
	Create(ctx context.Context, customer model.Customer) error
	Update(ctx context.Context, customer model.Customer) error
}
