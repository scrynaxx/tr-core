package middleware

import (
	"uuid"

	"github.com/labstack/echo/v5"
)

const HeaderXRequestID = "X-Request-ID"

// RequestID добавляет request ID к входящему запросу.
//
// Если клиент передал X-Request-ID, используется его значение.
// В противном случае генерируется новый UUID.
//
// Request ID устанавливается одновременно в request и response headers,
// чтобы быть доступным внутри приложения и возвращаться клиенту.
func RequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			requestID := c.Request().Header.Get(HeaderXRequestID)
			if requestID == "" {
				requestID = uuid.New().String()
			}

			c.Request().Header.Set(HeaderXRequestID, requestID)
			c.Response().Header().Set(HeaderXRequestID, requestID)

			return next(c)
		}
	}
}
