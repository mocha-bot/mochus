package http_middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	zLog "github.com/rs/zerolog/log"
)

func RequestLoggerWithZerolog() middleware.RequestLoggerConfig {
	return middleware.RequestLoggerConfig{
		LogURI:           true,
		LogStatus:        true,
		LogMethod:        true,
		LogUserAgent:     true,
		LogLatency:       true,
		LogError:         true,
		LogRemoteIP:      true,
		LogProtocol:      true,
		LogHost:          true,
		LogRequestID:     true,
		LogReferer:       true,
		LogContentLength: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			zLog.Info().
				Str("method", v.Method).
				Int("status", v.Status).
				Str("host", v.Host).
				Str("uri", v.URI).
				Str("referer", v.Referer).
				Str("user_agent", v.UserAgent).
				Dur("latency", v.Latency).
				Str("remote_ip", v.RemoteIP).
				Str("protocol", v.Protocol).
				Str("request_id", v.RequestID).
				Str("content_length", v.ContentLength).
				Msg("Request received")
			return nil
		},
	}
}
