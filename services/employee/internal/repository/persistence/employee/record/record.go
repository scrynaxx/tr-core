package record

import (
	"time"

	"github.com/google/uuid"
	"github.com/scrynaxx/tr-core/services/employee/internal/model"
)

type Employee struct {
	ID         uuid.UUID          `db:"id"`
	Type       model.EmployeeType `db:"type"`
	FirstName  string             `db:"first_name"`
	LastName   string             `db:"last_name"`
	Patronymic string             `db:"patronymic"`
	Phone      string             `db:"phone"`
	BirthDate  time.Time          `db:"birth_date"`
	CreatedAt  time.Time          `db:"created_at"`
	UpdatedAt  time.Time          `db:"updated_at"`
	ArchivedAt *time.Time         `db:"archived_at"`
}
