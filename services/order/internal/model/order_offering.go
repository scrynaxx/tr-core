package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// OrderOffering Услуга привязанная к заказу.
type OrderOffering struct {
	OfferingID        *uuid.UUID      `json:"offering_id"`
	OfferingVariantID *uuid.UUID      `json:"offering_variant_id"`
	Name              string          `json:"name"`
	Quantity          int             `json:"quantity"`
	UnitPrice         decimal.Decimal `json:"unit_price"`
}
