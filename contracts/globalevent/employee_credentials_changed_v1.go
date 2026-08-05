package globalevent

import (
	"github.com/google/uuid"
	"github.com/scrynaxx/tr-core/pkg/events"
)

// EmployeeCredentialsChangedV1 глобальное событие изменения учётных данных сотрудника.
var EmployeeCredentialsChangedV1 = events.NewDescriptor[EmployeeCredentialsChangedDataV1]("employee.employee.credentials.changed.v1")

type EmployeeCredentialsChangedDataV1 struct {
	EmployeeID uuid.UUID `json:"employee_id"`
}
