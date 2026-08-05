package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// PackingMaterial Упаковка / упаковочный материал.
type PackingMaterial struct {
	ID         uuid.UUID       `json:"id"`
	Name       string          `json:"name"`
	Price      decimal.Decimal `json:"price"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
	ArchivedAt *time.Time      `json:"archived_at"`
}
