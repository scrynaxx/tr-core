package record

import "github.com/scrynaxx/tr-core/services/customer/internal/model"

func ToCustomer(rec Customer) model.Customer {
	return model.Customer{
		ID:         rec.ID,
		Name:       rec.Name,
		Phone:      rec.Phone,
		Email:      rec.Email,
		CreatedAt:  rec.CreatedAt,
		UpdatedAt:  rec.UpdatedAt,
		ArchivedAt: rec.ArchivedAt,
	}
}

func ToCustomers(recs []Customer) []model.Customer {
	customers := make([]model.Customer, len(recs))
	for i := range recs {
		customers[i] = ToCustomer(recs[i])
	}

	return customers
}
