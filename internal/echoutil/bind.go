package echoutil

import (
	"fmt"

	"github.com/labstack/echo/v5"
)

func Bind[T any](c *echo.Context) (T, error) {
	var req T
	if err := c.Bind(&req); err != nil {
		var zero T
		return zero, fmt.Errorf("failed to bind request: %w", err)
	}

	if err := c.Validate(&req); err != nil {
		var zero T
		return zero, fmt.Errorf("falied to validate request: %w", err)
	}

	return req, nil
}
