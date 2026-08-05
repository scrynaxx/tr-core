package model

import (
	"uuid"

	"github.com/shopspring/decimal"
)

// OrderOffering Услуга привязанная к заказу.
type OrderOffering struct {
	OfferingID *uuid.UUID      `json:"offering_id"`
	Name       string          `json:"name"`
	Quantity   int             `json:"quantity"`
	UnitPrice  decimal.Decimal `json:"unit_price"`
}
