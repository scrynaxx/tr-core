package record

import (
	"time"

	"github.com/google/uuid"
)

type Passport struct {
	EmployeeID     uuid.UUID `db:"employee_id"`
	Series         string    `db:"series"`
	Number         string    `db:"number"`
	IssuedBy       string    `db:"issued_by"`
	IssuedAt       time.Time `db:"issued_at"`
	DepartmentCode string    `db:"department_code"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}
