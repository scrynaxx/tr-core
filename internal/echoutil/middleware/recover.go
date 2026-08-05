package middleware

import (
	"fmt"
	"log/slog"
	"runtime/debug"

	"github.com/labstack/echo/v5"
)

// Recover возвращает Echo middleware, который перехватывает panic,
// записывает информацию о нём в лог вместе со stack trace и преобразует
// panic во внутреннюю ошибку сервера.
func Recover() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) (err error) {
			defer func() {
				if r := recover(); r != nil {
					slog.Error("[recover middleware] panic recovered",
						slog.Any("panic", r),
						slog.String("stack", string(debug.Stack())),
					)

					err = fmt.Errorf("panic recovered: %v", r)
				}
			}()

			return next(c)
		}
	}
}
