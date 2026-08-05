package rest

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/scrynaxx/tr-core/config"
	"github.com/scrynaxx/tr-core/internal/echoutil/middleware"
	"github.com/scrynaxx/tr-core/internal/modules/auth/controller/rest/v1"
	"github.com/scrynaxx/tr-core/internal/modules/auth/model"
	"github.com/scrynaxx/tr-core/internal/modules/auth/usecase"
	"github.com/scrynaxx/tr-core/internal/security"
)

func RegisterRoutes(e *echo.Echo, authorizer security.Authorizer, env config.Environment, authUc usecase.Auth) {
	e.Use(middleware.ErrorMapper(map[error]int{
		model.ErrInvalidCredentials: http.StatusBadRequest,
		model.ErrSessionExpired:     http.StatusUnauthorized,
		model.ErrSessionNotFound:    http.StatusUnauthorized,
		model.ErrEmailAlreadyTaken:  http.StatusConflict,
		model.ErrAccountNotFound:    http.StatusUnauthorized,
	}))

	c := v1.New(env, authUc)
	adminGroupV1 := e.Group("/v1/admin/auth")
	webGroupV1 := e.Group("/v1/web/auth")

	adminGroupV1.POST("/sign-in", c.SignIn)
	adminGroupV1.POST("/refresh", c.Refresh)
	adminGroupV1.POST("/sign-out", c.SignOut, middleware.Auth(authorizer, security.ActorEmployee))

	webGroupV1.POST("/sign-in", c.SignIn)
	webGroupV1.POST("/refresh", c.Refresh)
	webGroupV1.POST("/sign-out", c.SignOut, middleware.Auth(authorizer, security.ActorCustomer))
}
