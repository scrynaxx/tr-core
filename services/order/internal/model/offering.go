package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Offering Дополнительная услуга к заказу
type Offering struct {
	ID         uuid.UUID        `json:"id"`
	Name       string           `json:"name"`
	Price      *decimal.Decimal `json:"price"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
	ArchivedAt *time.Time       `json:"archived_at"`
}
