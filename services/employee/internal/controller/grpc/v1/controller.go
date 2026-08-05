package v1

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	pbv1 "github.com/scrynaxx/tr-core/contracts/generated/services/v1/employee"
	"github.com/scrynaxx/tr-core/pkg/transport"
	"github.com/scrynaxx/tr-core/services/employee/internal/model"
	"github.com/scrynaxx/tr-core/services/employee/internal/usecase"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Controller struct {
	employeeCase usecase.Employee
}

func NewController(employeeCase usecase.Employee) pbv1.EmployeeServiceServer {
	return &Controller{employeeCase: employeeCase}
}

func (c *Controller) GetMe(ctx context.Context, empty *emptypb.Empty) (*pbv1.Employee, error) {
	identity, err := transport.IdentityFromContext(ctx)
	if err != nil {
		return nil, err
	}

	employee, err := c.employeeCase.Get(ctx, identity.EmployeeID)
	if err != nil {
		return nil, err
	}

	return toEmployee(employee), nil
}

func (c *Controller) CreateEmployee(ctx context.Context, req *pbv1.CreateEmployeeRequest) (*pbv1.Employee, error) {
	params := model.CreateEmployeeParams{
		Type:        model.EmployeeType(req.Type),
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Patronymic:  req.Patronymic,
		Phone:       req.Phone,
		BirthDate:   req.BirthDate.AsTime(),
		Credentials: nil,
		Passport:    nil,
	}
	if req.Credentials != nil {
		params.Credentials = &model.CreateCredentialsParams{
			Email:    req.Credentials.Email,
			Password: req.Credentials.Password,
		}
	}
	if req.Passport != nil {
		params.Passport = &model.CreatePassportParams{
			Series:   req.Passport.Series,
			Number:   req.Passport.Number,
			IssuedBy: req.Passport.IssuedBy,
			IssuedAt: req.Passport.IssuedAt.AsTime(),
		}
	}

	employee, err := c.employeeCase.Create(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("create employee: %w", err)
	}

	return toEmployee(employee), nil
}

func (c *Controller) GetEmployee(ctx context.Context, req *pbv1.GetEmployeeRequest) (*pbv1.Employee, error) {
	employee, err := c.employeeCase.Get(ctx, uuid.MustParse(req.EmployeeId))
	if err != nil {
		return nil, fmt.Errorf("get employee: %w", err)
	}

	return toEmployee(employee), nil
}

func (c *Controller) ListEmployees(ctx context.Context, _ *emptypb.Empty) (*pbv1.ListEmployeesResponse, error) {
	employees, err := c.employeeCase.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list employees: %w", err)
	}

	return toListEmployeesResponse(employees), nil
}

func (c *Controller) UpdateEmployee(ctx context.Context, req *pbv1.UpdateEmployeeRequest) (*pbv1.Employee, error) {
	params := model.UpdateEmployeeParams{
		Type:       model.EmployeeType(req.Type),
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Patronymic: req.Patronymic,
		Phone:      req.Phone,
		BirthDate:  req.BirthDate.AsTime(),
	}

	employee, err := c.employeeCase.Update(ctx, uuid.MustParse(req.EmployeeId), params)
	if err != nil {
		return nil, fmt.Errorf("update employee: %w", err)
	}

	return toEmployee(employee), nil
}

func (c *Controller) ArchiveEmployee(ctx context.Context, req *pbv1.ArchiveEmployeeRequest) (*emptypb.Empty, error) {
	identity, err := transport.IdentityFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("identity from context: %w", err)
	}

	if err = c.employeeCase.Archive(ctx, identity.EmployeeID, uuid.MustParse(req.EmployeeId)); err != nil {
		return nil, fmt.Errorf("archive employee: %w", err)
	}

	return new(emptypb.Empty), nil
}

func (c *Controller) CreateCredentials(ctx context.Context, req *pbv1.CreateCredentialsRequest) (*pbv1.Credentials, error) {
	params := model.CreateCredentialsParams{
		Email:    req.Email,
		Password: req.Password,
	}

	credentials, err := c.employeeCase.CreateCredentials(ctx, uuid.MustParse(req.EmployeeId), params)
	if err != nil {
		return nil, fmt.Errorf("create credentials: %w", err)
	}

	return toCredentials(credentials), nil
}

func (c *Controller) FindCredentials(ctx context.Context, req *pbv1.FindCredentialsRequest) (*pbv1.Credentials, error) {
	credentials, err := c.employeeCase.FindCredentials(ctx, uuid.MustParse(req.EmployeeId))
	if err != nil {
		return nil, fmt.Errorf("find credentials: %w", err)
	}

	if credentials == nil {
		return nil, nil
	}

	return toCredentials(*credentials), nil
}

func (c *Controller) UpdateCredentials(ctx context.Context, req *pbv1.UpdateCredentialsRequest) (*pbv1.Credentials, error) {
	credentials, err := c.employeeCase.UpdateCredentials(ctx, uuid.MustParse(req.EmployeeId), req.Email)
	if err != nil {
		return nil, fmt.Errorf("create credentials: %w", err)
	}

	return toCredentials(credentials), nil
}

func (c *Controller) DeleteCredentials(ctx context.Context, req *pbv1.DeleteCredentialsRequest) (*emptypb.Empty, error) {
	if err := c.employeeCase.DeleteCredentials(ctx, uuid.MustParse(req.EmployeeId)); err != nil {
		return nil, fmt.Errorf("delete employee: %w", err)
	}

	return new(emptypb.Empty), nil
}

func (c *Controller) CreatePassport(ctx context.Context, req *pbv1.CreatePassportRequest) (*pbv1.Passport, error) {
	params := model.CreatePassportParams{
		Series:         req.Series,
		Number:         req.Number,
		IssuedBy:       req.IssuedBy,
		IssuedAt:       req.IssuedAt.AsTime(),
		DepartmentCode: req.DepartmentCode,
	}

	passport, err := c.employeeCase.CreatePassport(ctx, uuid.MustParse(req.EmployeeId), params)
	if err != nil {
		return nil, fmt.Errorf("create passport: %w", err)
	}

	return toPassport(passport), nil
}

func (c *Controller) FindPassport(ctx context.Context, req *pbv1.FindPassportRequest) (*pbv1.Passport, error) {
	passport, err := c.employeeCase.FindPassport(ctx, uuid.MustParse(req.EmployeeId))
	if err != nil {
		return nil, fmt.Errorf("find passport: %w", err)
	}

	if passport == nil {
		return nil, nil
	}

	return toPassport(*passport), nil
}

func (c *Controller) UpdatePassport(ctx context.Context, req *pbv1.UpdatePassportRequest) (*pbv1.Passport, error) {
	params := model.UpdatePassportParams{
		Series:         req.Series,
		Number:         req.Number,
		IssuedBy:       req.IssuedBy,
		IssuedAt:       req.IssuedAt.AsTime(),
		DepartmentCode: req.DepartmentCode,
	}

	passport, err := c.employeeCase.UpdatePassport(ctx, uuid.MustParse(req.EmployeeId), params)
	if err != nil {
		return nil, fmt.Errorf("update passport: %w", err)
	}

	return toPassport(passport), nil
}

func (c *Controller) DeletePassport(ctx context.Context, req *pbv1.DeletePassportRequest) (*emptypb.Empty, error) {
	if err := c.employeeCase.DeletePassport(ctx, uuid.MustParse(req.EmployeeId)); err != nil {
		return nil, fmt.Errorf("delete employee: %w", err)
	}

	return new(emptypb.Empty), nil
}
