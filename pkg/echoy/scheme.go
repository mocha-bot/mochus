package echoy

import "github.com/labstack/echo/v4"

const (
	HTTP  = "http"
	HTTPS = "https"
)

func GetScheme(c echo.Context) string {
	if c.Request().TLS != nil {
		return HTTPS
	}
	return HTTP
}
