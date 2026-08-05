package rest

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/scrynaxx/tr-core/internal/echoutil/middleware"
	authclient "github.com/scrynaxx/tr-core/internal/modules/auth/client"
	"github.com/scrynaxx/tr-core/internal/modules/employee/controller/rest/v1"
	"github.com/scrynaxx/tr-core/internal/modules/employee/model"
	"github.com/scrynaxx/tr-core/internal/modules/employee/usecase"
	"github.com/scrynaxx/tr-core/internal/security"
)

func RegisterRoutes(e *echo.Echo, authorizer security.Authorizer, employeeUc usecase.Employee, passportUc usecase.Passport, authClient authclient.Client) {
	e.Use(middleware.ErrorMapper(map[error]int{
		model.ErrEmployeeNotFound: http.StatusNotFound,
		model.ErrPassportNotFound: http.StatusNotFound,
	}))

	c := v1.New(employeeUc, passportUc, authClient)
	adminGroupV1 := e.Group("/v1/admin/employee", middleware.Auth(authorizer, security.ActorEmployee))

	adminGroupV1.GET("/profile", c.GetProfile)
	adminGroupV1.POST("/employees", c.Create)
	adminGroupV1.GET("/employees/:employee_id", c.Get)
	adminGroupV1.GET("/employees", c.List)
	adminGroupV1.PUT("/employees/:employee_id", c.Update)
	adminGroupV1.POST("/employees/:employee_id/archive", c.Archive)
	adminGroupV1.POST("/employees/:employee_id/restore", c.Restore)
	adminGroupV1.POST("/employees/:employee_id/passport", c.CreatePassport)
	adminGroupV1.GET("/employees/:employee_id/passport", c.FindPassport)
	adminGroupV1.PUT("/employees/:employee_id/passport", c.UpdatePassport)
	adminGroupV1.DELETE("/employees/:employee_id/passport", c.DeletePassport)
}
