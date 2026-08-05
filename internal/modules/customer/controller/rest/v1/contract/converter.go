package contract

import (
	"github.com/scrynaxx/tr-core/internal/modules/customer/model"
)

func ToCustomer(customer model.Customer) Customer {
	return Customer{
		ID:         customer.ID,
		AccountID:  customer.AccountID,
		FirstName:  customer.FirstName,
		LastName:   customer.LastName,
		Patronymic: customer.Patronymic,
		Phone:      customer.Phone,
		Email:      customer.Email,
		CreatedAt:  customer.CreatedAt,
		UpdatedAt:  customer.UpdatedAt,
		ArchivedAt: customer.ArchivedAt,
	}
}

func ToListCustomerResponse(customers []model.Customer) ListCustomersResponse {
	res := ListCustomersResponse{
		Customers: make([]Customer, len(customers)),
	}

	for i, customer := range customers {
		res.Customers[i] = ToCustomer(customer)
	}

	return res
}
