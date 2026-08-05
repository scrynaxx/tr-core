package inner

import (
	"context"

	pbinner "github.com/scrynaxx/tr-core/contracts/generated/services/inner/employee"
	"github.com/scrynaxx/tr-core/services/employee/internal/usecase"
)

type Controller struct {
	employeeCase usecase.Employee
}

func NewController(employeeCase usecase.Employee) pbinner.EmployeeServiceServer {
	return &Controller{employeeCase: employeeCase}
}

func (c *Controller) FindCredentials(ctx context.Context, req *pbinner.FindCredentialsRequest) (*pbinner.Credentials, error) {
	employeeID, credentials, err := c.employeeCase.FindCredentialsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	return toCredentials(employeeID, credentials), nil
}
