package model

import "github.com/shopspring/decimal"

// PaymentPart Оплата заказа.
type PaymentPart struct {
	Amount decimal.Decimal `json:"amount"`
	Method PaymentMethod   `json:"method"`
}
