package middleware

import (
	"log/slog"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/scrynaxx/tr-core/internal/security"
)

// Auth возвращает Echo middleware для аутентификации и авторизации запросов с использованием Bearer token.
//
// Middleware проверяет token, статус соответствующей сессии через security.Authorizer и тип пользователя.
// Ошибки аутентификации и авторизации возвращаются без преобразования и
// обрабатываются middleware сопоставления ошибок, если он установлен.
//
// После успешной авторизации Identity сохраняется в Echo context.
func Auth(authorizer security.Authorizer, actor security.Actor) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			const prefix = "Bearer "

			value := c.Request().Header.Get(echo.HeaderAuthorization)
			if !strings.HasPrefix(value, prefix) {
				return security.ErrInvalidToken
			}

			token := strings.TrimSpace(strings.TrimPrefix(value, prefix))
			if token == "" {
				return security.ErrInvalidToken
			}

			identity, err := authorizer.Authorize(c.Request().Context(), token)
			if err != nil {
				return err
			}

			if identity.Actor != actor {
				slog.Info("[auth middleware] actor attempted to access a forbidden resource",
					slog.String("account_id", identity.AccountID.String()),
					slog.String("actor", string(identity.Actor)),
					slog.String("required_actor", string(actor)),
				)
				return security.ErrForbidden
			}

			security.SetIdentity(c, identity)

			return next(c)
		}
	}
}
