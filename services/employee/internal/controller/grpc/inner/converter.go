package inner

import (
	"github.com/google/uuid"
	pbinner "github.com/scrynaxx/tr-core/contracts/generated/services/inner/employee"
	"github.com/scrynaxx/tr-core/services/employee/internal/model"
)

func toCredentials(employeeID uuid.UUID, credentials *model.EmployeeCredentials) *pbinner.Credentials {
	if credentials == nil {
		return nil
	}

	return &pbinner.Credentials{
		EmployeeId:   employeeID.String(),
		Email:        credentials.Email,
		PasswordHash: credentials.PasswordHash,
	}
}
