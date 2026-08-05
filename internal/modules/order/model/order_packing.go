package model

import (
	"github.com/shopspring/decimal"
	"uuid"
)

// OrderPacking Упаковка привязанная к заказу.
type OrderPacking struct {
	PackingID *uuid.UUID      `json:"packing_id"`
	Name      string          `json:"name"`
	Quantity  int             `json:"quantity"`
	UnitPrice decimal.Decimal `json:"unit_price"`
}
