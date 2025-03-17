package http_middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type CORSOptions struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowCredentials bool
}

type CORSOption func(*CORSOptions)

func WithAllowOrigins(origins []string) CORSOption {
	return func(o *CORSOptions) {
		o.AllowOrigins = origins
	}
}

func WithAllowMethods(methods []string) CORSOption {
	return func(o *CORSOptions) {
		o.AllowMethods = methods
	}
}

func WithAllowCredentials(allow bool) CORSOption {
	return func(o *CORSOptions) {
		o.AllowCredentials = allow
	}
}

func CORS(opts ...CORSOption) echo.MiddlewareFunc {
	options := &CORSOptions{
		AllowOrigins:     []string{},
		AllowMethods:     []string{},
		AllowCredentials: false,
	}

	for _, opt := range opts {
		opt(options)
	}

	return middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:          middleware.DefaultSkipper,
		AllowOrigins:     options.AllowOrigins,
		AllowMethods:     options.AllowMethods,
		AllowCredentials: options.AllowCredentials,
	})
}
