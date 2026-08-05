package model

// PaymentMethod Тип оплаты.
type PaymentMethod string

const (
	// PaymentMethodWithVAT Оплата с НДС.
	PaymentMethodWithVAT PaymentMethod = "with_vat"

	// PaymentMethodWithoutVAT Оплата без НДС.
	PaymentMethodWithoutVAT PaymentMethod = "without_vat"

	// PaymentMethodCash Оплата наличными.
	PaymentMethodCash PaymentMethod = "cash"
)
