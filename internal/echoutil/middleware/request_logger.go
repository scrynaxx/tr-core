package middleware

import (
	"context"
	"log/slog"
	"time"

	"github.com/labstack/echo/v5"
)

// RequestLogger логирует каждый HTTP-запрос с информацией о методе,
// URI, статусе ответа, времени выполнения, host, user agent и real IP.
//
// Если запрос содержит X-Request-ID, его значение также добавляется в лог.
// Ошибки логируются с уровнем ErrorDetailer и возвращаются дальше по middleware chain для последующей обработки.
func RequestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			start := time.Now()
			err := next(c)

			status := -1
			if response := c.Request().Response; response != nil && response.StatusCode != 0 {
				status = response.StatusCode
			}

			attrs := []slog.Attr{
				slog.String("method", c.Request().Method),
				slog.String("uri", c.Request().RequestURI),
				slog.Int("status", status),
				slog.Int64("latency_ms", time.Since(start).Milliseconds()),
				slog.String("host", c.Request().Host),
				slog.String("user_agent", c.Request().UserAgent()),
				slog.String("real_ip", c.RealIP()),
			}

			requestID := c.Request().Header.Get(HeaderXRequestID)
			if requestID != "" {
				attrs = append(attrs, slog.String("request_id", requestID))
			}

			level := slog.LevelInfo
			if err != nil {
				level = slog.LevelError
				attrs = append(attrs, slog.String("error", err.Error()))
			}

			slog.LogAttrs(context.Background(), level, "[request logger middleware] request", attrs...)

			return err
		}
	}
}
