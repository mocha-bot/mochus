package http_handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mocha-bot/mochus/core/entity"
	cookiey "github.com/mocha-bot/mochus/pkg/cookie"
	zLog "github.com/rs/zerolog/log"
)

func parseOauthCallbackRequest(c echo.Context) (*entity.OauthCallbackRequest, error) {
	req := new(entity.OauthCallbackRequest)

	if err := c.Bind(req); err != nil {
		return nil, fmt.Errorf("%w: %s", entity.ErrorBind, err)
	}

	return req, nil
}

func parseOauthCallbackError(err error) (code int, i any) {
	switch {
	case errors.Is(err, entity.ErrorBind):
		return http.StatusBadRequest, Response{Message: err.Error()}
	default:
		zLog.Error().Err(err).Msg("Internal server error")
		return http.StatusInternalServerError, Response{Message: "Internal server error"}
	}
}

func parseRefreshTokenRequest(c echo.Context) (*entity.RefreshTokenRequest, error) {
	req := new(entity.RefreshTokenRequest)

	refreshTokenCookie, err := c.Cookie(cookiey.CookieRefreshToken)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", entity.ErrorUnauthorized, err)
	}

	req.RefreshToken = refreshTokenCookie.Value

	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, req); err != nil {
		return nil, fmt.Errorf("%w: %s", entity.ErrorUnauthorized, err)
	}

	return req, nil
}

func parseRefreshTokenError(err error) (code int, i any) {
	switch {
	case errors.Is(err, entity.ErrorUnauthorized):
		return http.StatusUnauthorized, Response{Message: err.Error()}
	default:
		zLog.Error().Err(err).Msg("Internal server error")
		return http.StatusInternalServerError, Response{Message: "Internal server error"}
	}
}

func parseRevokeTokenRequest(c echo.Context) (*entity.RevokeTokenRequest, error) {
	req := new(entity.RevokeTokenRequest)

	// error isn't checked because if cookie doesn't exist, it's still allowed
	refreshTokenCookie, _ := c.Cookie(cookiey.CookieRefreshToken)
	accessTokenCookie, _ := c.Cookie(cookiey.CookieAccessToken)

	req.RefreshToken = refreshTokenCookie.Value
	req.AccessToken = accessTokenCookie.Value

	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, req); err != nil {
		return nil, fmt.Errorf("%w: %s", entity.ErrorUnauthorized, err)
	}

	return req, nil
}

func parseRevokeTokenError(err error) (code int, i any) {
	switch {
	case errors.Is(err, entity.ErrorUnauthorized):
		return http.StatusUnauthorized, Response{Message: err.Error()}
	default:
		zLog.Error().Err(err).Msg("Internal server error")
		return http.StatusInternalServerError, Response{Message: "Internal server error"}
	}
}

func parseGetUserByTokenRequest(c echo.Context) (*entity.GetUserByTokenRequest, error) {
	req := new(entity.GetUserByTokenRequest)

	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, req); err != nil {
		return nil, fmt.Errorf("%w: %s", entity.ErrorUnauthorized, err)
	}

	return req, nil
}

func parseGetUserByTokenError(err error) (code int, i any) {
	switch {
	case errors.Is(err, entity.ErrorUnauthorized):
		return http.StatusUnauthorized, Response{Message: err.Error()}
	default:
		zLog.Error().Err(err).Msg("Internal server error")
		return http.StatusInternalServerError, Response{Message: "Internal server error"}
	}
}

func parseGetUserByTokenResponse(user *entity.User) (code int, i any) {
	return http.StatusOK, Response{Data: user, Message: "Success retrieve user"}
}
