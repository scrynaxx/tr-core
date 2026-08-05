package record

import (
	"github.com/scrynaxx/tr-core/internal/modules/employee/model"
)

func ToPassport(rec Passport) model.Passport {
	return model.Passport{
		Series:         rec.Series,
		Number:         rec.Number,
		IssuedBy:       rec.IssuedBy,
		IssuedAt:       rec.IssuedAt,
		DepartmentCode: rec.DepartmentCode,
		CreatedAt:      rec.CreatedAt,
		UpdatedAt:      rec.UpdatedAt,
	}
}

func ToPassportPtr(rec *Passport) *model.Passport {
	if rec == nil {
		return nil
	}

	return &model.Passport{
		Series:         rec.Series,
		Number:         rec.Number,
		IssuedBy:       rec.IssuedBy,
		IssuedAt:       rec.IssuedAt,
		DepartmentCode: rec.DepartmentCode,
		CreatedAt:      rec.CreatedAt,
		UpdatedAt:      rec.UpdatedAt,
	}
}
