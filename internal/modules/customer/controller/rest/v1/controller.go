package v1

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/scrynaxx/tr-core/internal/echoutil"
	"github.com/scrynaxx/tr-core/internal/modules/customer/controller/rest/v1/contract"
	"github.com/scrynaxx/tr-core/internal/modules/customer/model"
	"github.com/scrynaxx/tr-core/internal/modules/customer/usecase"
)

type Controller struct {
	customerUc usecase.Customer
}

func New(customerUc usecase.Customer) *Controller {
	return &Controller{
		customerUc: customerUc,
	}
}

func (co *Controller) Create(c *echo.Context) error {
	req, err := echoutil.Bind[contract.CreateCustomerRequest](c)
	if err != nil {
		return err
	}

	if err = co.customerUc.Create(c.Request().Context(), model.CustomerInput{
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Patronymic: req.Patronymic,
		Phone:      req.Phone,
		Email:      req.Email,
	}); err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)
}

func (co *Controller) Get(c *echo.Context) error {
	req, err := echoutil.Bind[contract.GetCustomerRequest](c)
	if err != nil {
		return err
	}

	employee, err := co.customerUc.Get(c.Request().Context(), req.CustomerID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, contract.ToCustomer(employee))
}

func (co *Controller) List(c *echo.Context) error {
	employees, err := co.customerUc.List(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, contract.ToListCustomerResponse(employees))
}

func (co *Controller) Update(c *echo.Context) error {
	req, err := echoutil.Bind[contract.UpdateCustomerRequest](c)
	if err != nil {
		return err
	}

	if err = co.customerUc.Update(c.Request().Context(), req.CustomerID, model.CustomerInput{
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Patronymic: req.Patronymic,
		Phone:      req.Phone,
		Email:      req.Email,
	}); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (co *Controller) Archive(c *echo.Context) error {
	req, err := echoutil.Bind[contract.ArchiveCustomerRequest](c)
	if err != nil {
		return err
	}

	if err = co.customerUc.Archive(c.Request().Context(), req.CustomerID); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (co *Controller) Restore(c *echo.Context) error {
	req, err := echoutil.Bind[contract.RestoreCustomerRequest](c)
	if err != nil {
		return err
	}

	if err = co.customerUc.Restore(c.Request().Context(), req.CustomerID); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
