package contract

import (
	authclient "github.com/scrynaxx/tr-core/internal/modules/auth/client"
	"github.com/scrynaxx/tr-core/internal/modules/employee/model"
)

func ToEmployee(employee model.Employee) Employee {
	return Employee{
		ID:         employee.ID,
		AccountID:  employee.AccountID,
		Type:       employee.Type,
		FirstName:  employee.FirstName,
		LastName:   employee.LastName,
		Patronymic: employee.Patronymic,
		Phone:      employee.Phone,
		BirthDate:  employee.BirthDate,
		CreatedAt:  employee.CreatedAt,
		UpdatedAt:  employee.UpdatedAt,
		ArchivedAt: employee.ArchivedAt,
	}
}

func ToListEmployeeResponse(employees []model.Employee) ListEmployeesResponse {
	res := ListEmployeesResponse{
		Employees: make([]Employee, len(employees)),
	}

	for i, employee := range employees {
		res.Employees[i] = ToEmployee(employee)
	}

	return res
}

func ToPassport(passport *model.Passport) *Passport {
	if passport == nil {
		return nil
	}

	return &Passport{
		Series:         passport.Series,
		Number:         passport.Number,
		IssuedBy:       passport.IssuedBy,
		IssuedAt:       passport.IssuedAt,
		DepartmentCode: passport.DepartmentCode,
		CreatedAt:      passport.CreatedAt,
		UpdatedAt:      passport.UpdatedAt,
	}
}

func ToProfile(employee model.Employee, account authclient.Account) Profile {
	return Profile{
		ID:         employee.ID,
		AccountID:  account.ID,
		Type:       employee.Type,
		FirstName:  employee.FirstName,
		LastName:   employee.LastName,
		Patronymic: employee.Patronymic,
		Phone:      employee.Phone,
		BirthDate:  employee.BirthDate,
		CreatedAt:  employee.CreatedAt,
		UpdatedAt:  employee.UpdatedAt,
		Email:      account.Email,
	}
}
