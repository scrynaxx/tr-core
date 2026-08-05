package record

import (
	"github.com/scrynaxx/tr-core/internal/modules/employee/model"
)

func ToEmployee(rec Employee) model.Employee {
	return model.Employee{
		ID:         rec.ID,
		AccountID:  rec.AccountID,
		Type:       rec.Type,
		FirstName:  rec.FirstName,
		LastName:   rec.LastName,
		Patronymic: rec.Patronymic,
		Phone:      rec.Phone,
		BirthDate:  rec.BirthDate,
		CreatedAt:  rec.CreatedAt,
		UpdatedAt:  rec.UpdatedAt,
		ArchivedAt: rec.ArchivedAt,
	}
}

func ToListEmployee(recs []Employee) []model.Employee {
	list := make([]model.Employee, len(recs))
	for i, r := range recs {
		list[i] = ToEmployee(r)
	}

	return list
}
