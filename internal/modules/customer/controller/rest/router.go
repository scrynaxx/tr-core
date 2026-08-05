package rest

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/scrynaxx/tr-core/internal/echoutil/middleware"
	v1 "github.com/scrynaxx/tr-core/internal/modules/customer/controller/rest/v1"
	"github.com/scrynaxx/tr-core/internal/modules/customer/model"
	"github.com/scrynaxx/tr-core/internal/modules/customer/usecase"
	"github.com/scrynaxx/tr-core/internal/security"
)

func RegisterRoutes(e *echo.Echo, authorizer security.Authorizer, customerUc usecase.Customer) {
	e.Use(middleware.ErrorMapper(map[error]int{
		model.ErrCustomerNotFound:      http.StatusNotFound,
		model.ErrCustomerAlreadyExists: http.StatusConflict,
	}))

	c := v1.New(customerUc)
	adminGroupV1 := e.Group("/v1/admin/customer", middleware.Auth(authorizer, security.ActorEmployee))

	adminGroupV1.POST("/customers", c.Create)
	adminGroupV1.GET("/customers/:customer_id", c.Get)
	adminGroupV1.GET("/customers", c.List)
	adminGroupV1.PUT("/customers/:customer_id", c.Update)
	adminGroupV1.POST("/customers/:customer_id/archive", c.Archive)
	adminGroupV1.POST("/customers/:customer_id/restore", c.Restore)
}
