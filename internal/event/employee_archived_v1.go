package event

import (
	"uuid"

	"github.com/scrynaxx/tr-core/pkg/messaging"
)

var (
	EmployeeArchivedV1 = messaging.NewDescriptor[EmployeeArchivedV1Payload]("employee.employee.archived.v1")
)

type EmployeeArchivedV1Payload struct {
	AccountID  uuid.UUID `json:"account_id"`
	EmployeeID uuid.UUID `json:"employee_id"`
}
