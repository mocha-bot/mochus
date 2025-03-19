package http_middleware

import (
	"github.com/labstack/echo/v4"
)

func FallbackRedirect(host string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Request().Header.Set("X-Fallback-Host", host)
			return next(c)
		}
	}
}
