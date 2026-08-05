package record

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID         uuid.UUID  `db:"id"`
	Name       string     `db:"name"`
	Phone      string     `db:"phone"`
	Email      string     `db:"email"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
	ArchivedAt *time.Time `db:"archived_at"`
}
