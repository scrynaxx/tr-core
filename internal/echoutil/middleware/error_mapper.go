package middleware

import (
	"errors"

	"github.com/labstack/echo/v5"
)

// ErrorMapper преобразует sentinel ошибки приложения в echo.HTTPError.
//
// Middleware не формирует HTTP-ответ самостоятельно, а только оборачивает
// ошибку из map в echo.HTTPError, если она найдена, и передаёт её дальше по цепочке middleware.
func ErrorMapper(mappings map[error]int) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			err := next(c)
			if err == nil {
				return nil
			}

			for target, status := range mappings {
				if errors.Is(err, target) {
					return echo.NewHTTPError(status, err.Error())
				}
			}

			return err
		}
	}
}
