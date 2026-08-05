package record

import (
	"github.com/scrynaxx/tr-core/services/employee/internal/model"
)

func ToCredentials(rec Credentials) model.EmployeeCredentials {
	return model.EmployeeCredentials{
		Email:        rec.Email,
		PasswordHash: rec.PasswordHash,
		CreatedAt:    rec.CreatedAt,
		UpdatedAt:    rec.UpdatedAt,
	}
}

func ToCredentialsPtr(rec *Credentials) *model.EmployeeCredentials {
	if rec == nil {
		return nil
	}

	return &model.EmployeeCredentials{
		Email:        rec.Email,
		PasswordHash: rec.PasswordHash,
		CreatedAt:    rec.CreatedAt,
		UpdatedAt:    rec.UpdatedAt,
	}
}
