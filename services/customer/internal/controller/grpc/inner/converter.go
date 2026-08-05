package inner

import (
	pbinner "github.com/scrynaxx/tr-core/contracts/generated/services/inner/customer"
	"github.com/scrynaxx/tr-core/services/customer/internal/model"
)

func toGetCustomerResponse(customer model.Customer) *pbinner.GetCustomerResponse {
	return &pbinner.GetCustomerResponse{
		Id:       customer.ID.String(),
		Name:     customer.Name,
		Phone:    customer.Phone,
		Email:    customer.Email,
		Archived: customer.ArchivedAt != nil,
	}
}
