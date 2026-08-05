package usecase

import (
	"context"

	"uuid"

	"github.com/scrynaxx/tr-core/internal/modules/customer/model"
)

type Customer interface {
	Create(ctx context.Context, input model.CustomerInput) error
	Get(ctx context.Context, customerID uuid.UUID) (model.Customer, error)
	List(ctx context.Context) ([]model.Customer, error)
	Update(ctx context.Context, customerID uuid.UUID, input model.CustomerInput) error
	Archive(ctx context.Context, customerID uuid.UUID) error
	Restore(ctx context.Context, customerID uuid.UUID) error
}
