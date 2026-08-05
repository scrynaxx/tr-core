package record

import (
	"github.com/scrynaxx/tr-core/services/employee/internal/model"
)

func ToEmployee(rec Employee) model.Employee {
	return model.Employee{
		ID:         rec.ID,
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

func ToEmployees(recs []Employee) []model.Employee {
	employees := make([]model.Employee, len(recs))
	for i := range recs {
		employees[i] = ToEmployee(recs[i])
	}

	return employees
}
