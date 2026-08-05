package globalevent

import (
	"github.com/google/uuid"
	"github.com/scrynaxx/tr-core/pkg/events"
)

// EmployeeArchivedV1 глобальное событие архивации сотрудника.
var EmployeeArchivedV1 = events.NewDescriptor[EmployeeArchivedDataV1]("employee.employee.archived.v1")

type EmployeeArchivedDataV1 struct {
	EmployeeID uuid.UUID `json:"employee_id"`
}
