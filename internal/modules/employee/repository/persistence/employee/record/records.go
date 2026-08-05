package record

import (
	"time"
	"uuid"

	"github.com/scrynaxx/tr-core/internal/modules/employee/model"
)

type Employee struct {
	ID         uuid.UUID          `db:"id"`
	AccountID  *uuid.UUID         `db:"account_id"`
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
