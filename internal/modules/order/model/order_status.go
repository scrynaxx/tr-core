package model

// OrderStatus Текущий статус заказа.
type OrderStatus string

const (
	// OrderStatusNew Новый заказ.
	OrderStatusNew OrderStatus = "new"

	// OrderStatusInProgress Заказ в работе.
	OrderStatusInProgress OrderStatus = "in_progress"

	// OrderStatusConfirmed Заказ подтвержден менеджером.
	OrderStatusConfirmed OrderStatus = "confirmed"

	// OrderStatusCancelled Заказ отменён.
	OrderStatusCancelled OrderStatus = "cancelled"

	// OrderStatusCompleted Заказ полностью выполнен.
	OrderStatusCompleted OrderStatus = "completed"
)
