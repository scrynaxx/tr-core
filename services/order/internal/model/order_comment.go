package model

import (
	"time"

	"github.com/google/uuid"
)

// OrderComment Комментарий к заказу.
type OrderComment struct {
	ID         uuid.UUID `json:"id"`
	EmployeeID uuid.UUID `json:"employee_id"`
	Value      string    `json:"value"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
