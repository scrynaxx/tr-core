package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/scrynaxx/tr-core/services/customer/internal/model"
)

// Customer описывает хранение клиентов.
type Customer interface {
	// Get возвращает клиента по идентификатору, включая архивного.
	Get(ctx context.Context, customerID uuid.UUID) (model.Customer, error)

	// List возвращает активных клиентов.
	List(ctx context.Context) ([]model.Customer, error)

	// Create сохраняет нового клиента.
	Create(ctx context.Context, customer model.Customer) (model.Customer, error)

	// Save изменяет активного клиента.
	Save(ctx context.Context, customer model.Customer) (model.Customer, error)

	// Archive исключает клиента из активного списка.
	Archive(ctx context.Context, customerID uuid.UUID) error
}
