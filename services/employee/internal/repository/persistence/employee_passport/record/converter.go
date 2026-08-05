package record

import (
	"github.com/scrynaxx/tr-core/services/employee/internal/model"
)

func ToPassport(rec Passport) model.EmployeePassport {
	return model.EmployeePassport{
		Series:         rec.Series,
		Number:         rec.Number,
		IssuedBy:       rec.IssuedBy,
		IssuedAt:       rec.IssuedAt,
		DepartmentCode: rec.DepartmentCode,
	}
}

func ToPassportPtr(rec *Passport) *model.EmployeePassport {
	if rec == nil {
		return nil
	}

	return &model.EmployeePassport{
		Series:         rec.Series,
		Number:         rec.Number,
		IssuedBy:       rec.IssuedBy,
		IssuedAt:       rec.IssuedAt,
		DepartmentCode: rec.DepartmentCode,
	}
}
