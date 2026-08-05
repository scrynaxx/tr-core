package v1

import (
	pbv1 "github.com/scrynaxx/tr-core/contracts/generated/services/v1/employee"
	"github.com/scrynaxx/tr-core/services/employee/internal/model"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func toEmployee(employee model.Employee) *pbv1.Employee {
	result := &pbv1.Employee{
		Id:         employee.ID.String(),
		Type:       string(employee.Type),
		FirstName:  employee.FirstName,
		LastName:   employee.LastName,
		Patronymic: employee.Patronymic,
		Phone:      employee.Phone,
		BirthDate:  timestamppb.New(employee.BirthDate),
		CreatedAt:  timestamppb.New(employee.CreatedAt),
		UpdatedAt:  timestamppb.New(employee.UpdatedAt),
		ArchivedAt: nil,
	}
	if employee.ArchivedAt != nil {
		result.ArchivedAt = timestamppb.New(*employee.ArchivedAt)
	}

	return result
}

func toListEmployeesResponse(employees []model.Employee) *pbv1.ListEmployeesResponse {
	result := &pbv1.ListEmployeesResponse{
		Employees: make([]*pbv1.Employee, len(employees)),
	}

	for i := range employees {
		result.Employees[i] = toEmployee(employees[i])
	}

	return result
}

func toCredentials(credentials model.EmployeeCredentials) *pbv1.Credentials {
	return &pbv1.Credentials{
		Email: credentials.Email,
	}
}

func toPassport(passport model.EmployeePassport) *pbv1.Passport {
	return &pbv1.Passport{
		Series:         passport.Series,
		Number:         passport.Number,
		IssuedBy:       passport.IssuedBy,
		IssuedAt:       timestamppb.New(passport.IssuedAt),
		DepartmentCode: passport.DepartmentCode,
	}
}
