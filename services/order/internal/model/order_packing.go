package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// OrderPacking Упаковка привязанная к заказу.
type OrderPacking struct {
	PackingID *uuid.UUID      `json:"packing_id"`
	Name      string          `json:"name"`
	Quantity  int             `json:"quantity"`
	UnitPrice decimal.Decimal `json:"unit_price"`
}
