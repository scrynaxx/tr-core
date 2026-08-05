package record

import "github.com/scrynaxx/tr-core/internal/modules/customer/model"

func ToCustomer(rec Customer) model.Customer {
	return model.Customer{
		ID:         rec.ID,
		AccountID:  rec.AccountID,
		FirstName:  rec.FirstName,
		LastName:   rec.LastName,
		Patronymic: rec.Patronymic,
		Phone:      rec.Phone,
		Email:      rec.Email,
		CreatedAt:  rec.CreatedAt,
		UpdatedAt:  rec.UpdatedAt,
		ArchivedAt: rec.ArchivedAt,
	}
}

func ToListCustomer(recs []Customer) []model.Customer {
	customers := make([]model.Customer, len(recs))
	for i := range recs {
		customers[i] = ToCustomer(recs[i])
	}

	return customers
}
