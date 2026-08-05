package v1

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/scrynaxx/tr-core/internal/echoutil"
	"github.com/scrynaxx/tr-core/internal/modules/auth/client"
	"github.com/scrynaxx/tr-core/internal/modules/employee/controller/rest/v1/contract"
	"github.com/scrynaxx/tr-core/internal/modules/employee/model"
	"github.com/scrynaxx/tr-core/internal/modules/employee/usecase"
	"github.com/scrynaxx/tr-core/internal/security"
)

type Controller struct {
	employeeUc usecase.Employee
	passportUc usecase.Passport
	authClient authclient.Client
}

func New(employeeUc usecase.Employee, passportUc usecase.Passport, authClient authclient.Client) *Controller {
	return &Controller{
		employeeUc: employeeUc,
		passportUc: passportUc,
		authClient: authClient,
	}
}

func (co *Controller) GetProfile(c *echo.Context) error {
	identity, err := security.GetIdentity(c)
	if err != nil {
		return err
	}

	employee, err := co.employeeUc.GetByAccount(c.Request().Context(), identity.AccountID)
	if err != nil {
		return err
	}

	account, err := co.authClient.GetAccount(c.Request().Context(), identity.AccountID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, contract.ToProfile(employee, account))
}

func (co *Controller) Create(c *echo.Context) error {
	req, err := echoutil.Bind[contract.CreateEmployeeRequest](c)
	if err != nil {
		return err
	}

	var passportInput *model.PassportInput
	if req.Passport != nil {
		passportInput = &model.PassportInput{
			Series:         req.Passport.Series,
			Number:         req.Passport.Number,
			IssuedBy:       req.Passport.IssuedBy,
			IssuedAt:       req.Passport.IssuedAt,
			DepartmentCode: req.Passport.DepartmentCode,
		}
	}

	if err = co.employeeUc.Create(c.Request().Context(), model.EmployeeInput{
		Type:       req.Type,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Patronymic: req.Patronymic,
		Phone:      req.Phone,
		BirthDate:  req.BirthDate,
	}, passportInput); err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)
}

func (co *Controller) Get(c *echo.Context) error {
	req, err := echoutil.Bind[contract.GetEmployeeRequest](c)
	if err != nil {
		return err
	}

	employee, err := co.employeeUc.Get(c.Request().Context(), req.EmployeeID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, contract.ToEmployee(employee))
}

func (co *Controller) List(c *echo.Context) error {
	employees, err := co.employeeUc.List(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, contract.ToListEmployeeResponse(employees))
}

func (co *Controller) Update(c *echo.Context) error {
	req, err := echoutil.Bind[contract.UpdateEmployeeRequest](c)
	if err != nil {
		return err
	}

	if err = co.employeeUc.Update(c.Request().Context(), req.EmployeeID, model.EmployeeInput{
		Type:       req.Type,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Patronymic: req.Patronymic,
		Phone:      req.Phone,
		BirthDate:  req.BirthDate,
	}); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (co *Controller) Archive(c *echo.Context) error {
	req, err := echoutil.Bind[contract.ArchiveEmployeeRequest](c)
	if err != nil {
		return err
	}

	identity, err := security.GetIdentity(c)
	if err != nil {
		return err
	}

	if err = co.employeeUc.Archive(c.Request().Context(), identity.AccountID, req.EmployeeID); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (co *Controller) Restore(c *echo.Context) error {
	req, err := echoutil.Bind[contract.RestoreEmployeeRequest](c)
	if err != nil {
		return err
	}

	if err = co.employeeUc.Restore(c.Request().Context(), req.EmployeeID); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (co *Controller) CreatePassport(c *echo.Context) error {
	req, err := echoutil.Bind[contract.CreatePassportRequest](c)
	if err != nil {
		return err
	}

	if err = co.passportUc.Create(c.Request().Context(), req.EmployeeID, model.PassportInput{
		Series:         req.Series,
		Number:         req.Number,
		IssuedBy:       req.IssuedBy,
		IssuedAt:       req.IssuedAt,
		DepartmentCode: req.DepartmentCode,
	}); err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)
}

func (co *Controller) FindPassport(c *echo.Context) error {
	req, err := echoutil.Bind[contract.GetPassportRequest](c)
	if err != nil {
		return err
	}

	passport, err := co.passportUc.Find(c.Request().Context(), req.EmployeeID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, contract.ToPassport(passport))
}

func (co *Controller) UpdatePassport(c *echo.Context) error {
	req, err := echoutil.Bind[contract.UpdatePassportRequest](c)
	if err != nil {
		return err
	}

	if err = co.passportUc.Update(c.Request().Context(), req.EmployeeID, model.PassportInput{
		Series:         req.Series,
		Number:         req.Number,
		IssuedBy:       req.IssuedBy,
		IssuedAt:       req.IssuedAt,
		DepartmentCode: req.DepartmentCode,
	}); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (co *Controller) DeletePassport(c *echo.Context) error {
	req, err := echoutil.Bind[contract.DeletePassportRequest](c)
	if err != nil {
		return err
	}

	if err = co.passportUc.Delete(c.Request().Context(), req.EmployeeID); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
