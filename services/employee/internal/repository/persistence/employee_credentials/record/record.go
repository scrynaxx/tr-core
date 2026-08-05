package record

import (
	"time"

	"github.com/google/uuid"
)

type Credentials struct {
	EmployeeID   uuid.UUID `db:"employee_id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
