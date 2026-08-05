package v1

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/scrynaxx/tr-core/config"
	"github.com/scrynaxx/tr-core/internal/echoutil"
	"github.com/scrynaxx/tr-core/internal/modules/auth/controller/rest/v1/contract"
	"github.com/scrynaxx/tr-core/internal/modules/auth/model"
	"github.com/scrynaxx/tr-core/internal/modules/auth/usecase"
	"github.com/scrynaxx/tr-core/internal/security"
)

const refreshCookieName = "refresh"

type Controller struct {
	env    config.Environment
	authUc usecase.Auth
}

func New(env config.Environment, authUc usecase.Auth) *Controller {
	return &Controller{
		env:    env,
		authUc: authUc,
	}
}

func (co *Controller) SignIn(c *echo.Context) error {
	req, err := echoutil.Bind[contract.SignInRequest](c)
	if err != nil {
		co.clearRefresh(c)
		return err
	}

	data, err := co.authUc.SignIn(c.Request().Context(), security.ActorEmployee, req.Email, req.Password, c.Request().UserAgent())
	if err != nil {
		co.clearRefresh(c)
		return err
	}

	co.applyRefresh(c, data)
	return c.JSON(http.StatusOK, contract.ToAuthResponse(data))
}

func (co *Controller) Refresh(c *echo.Context) error {
	refresh, err := c.Cookie(refreshCookieName)
	if err != nil {
		co.clearRefresh(c)
		return err
	}

	data, err := co.authUc.Refresh(c.Request().Context(), security.ActorEmployee, refresh.Value, c.Request().UserAgent())
	if err != nil {
		co.clearRefresh(c)
		return err
	}

	co.applyRefresh(c, data)
	return c.JSON(http.StatusOK, contract.ToAuthResponse(data))
}

func (co *Controller) SignOut(c *echo.Context) error {
	identity, err := security.GetIdentity(c)
	if err != nil {
		return err
	}

	if err = co.authUc.SignOut(c.Request().Context(), identity.AccountID, identity.SessionID); err != nil {
		return err
	}

	co.clearRefresh(c)
	return c.NoContent(http.StatusNoContent)
}

func (co *Controller) applyRefresh(c *echo.Context, res *model.AuthState) {
	sameSite := http.SameSiteNoneMode
	if co.env == config.Production {
		sameSite = http.SameSiteStrictMode
	}

	http.SetCookie(c.Response(), &http.Cookie{
		Name:     refreshCookieName,
		Path:     "/",
		Value:    res.RefreshToken,
		Expires:  res.RefreshUntil,
		HttpOnly: true,
		Secure:   true,
		SameSite: sameSite,
	})
}

func (co *Controller) clearRefresh(c *echo.Context) {
	sameSite := http.SameSiteNoneMode
	if co.env == config.Production {
		sameSite = http.SameSiteStrictMode
	}

	http.SetCookie(c.Response(), &http.Cookie{
		Name:     refreshCookieName,
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: sameSite,
	})
}
