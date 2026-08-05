package record

import (
	"time"

	"uuid"
)

type Customer struct {
	ID         uuid.UUID  `db:"id"`
	AccountID  *uuid.UUID `db:"account_id"`
	FirstName  string     `db:"first_name"`
	LastName   string     `db:"last_name"`
	Patronymic *string    `db:"patronymic"`
	Phone      string     `db:"phone"`
	Email      string     `db:"email"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
	ArchivedAt *time.Time `db:"archived_at"`
}
