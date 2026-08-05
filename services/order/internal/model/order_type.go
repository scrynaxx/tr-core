package model

// OrderType Тип заказа.
type OrderType string

const (
	// OrderTypeOfficeRelocation Заказ на офисный переезд.
	OrderTypeOfficeRelocation OrderType = "office_relocation"

	// OrderTypeFreightTransport Заказ по грузоперевозкам.
	OrderTypeFreightTransport OrderType = "freight_transport"
)
