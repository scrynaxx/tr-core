package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/scrynaxx/tr-core/services/customer/internal/model"
)

// Customer управляет клиентами.
type Customer interface {
	// Create создаёт клиента.
	Create(ctx context.Context, params model.CreateCustomerParams) (model.Customer, error)

	// Get возвращает клиента по идентификатору, включая архивного.
	Get(ctx context.Context, customerID uuid.UUID) (model.Customer, error)

	// List возвращает активных клиентов.
	List(ctx context.Context) ([]model.Customer, error)

	// Update изменяет активного клиента.
	Update(ctx context.Context, customerID uuid.UUID, params model.UpdateCustomerParams) (model.Customer, error)

	// Archive исключает клиента из дальнейшего выбора.
	Archive(ctx context.Context, customerID uuid.UUID) error
}
