package model

// OrderPayment Оплата заказа.
type OrderPayment struct {
	// Сколько и как платит клиент.
	Client PaymentPart `json:"client"`

	// Сколько и как получает исполнитель.
	Employee PaymentPart `json:"employee"`
}
