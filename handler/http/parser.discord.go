package http_handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mocha-bot/mochus/core/entity"
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
		return http.StatusInternalServerError, Response{Message: "Internal server error"}
	}
}

func parseRefreshTokenRequest(c echo.Context) (*entity.RefreshTokenRequest, error) {
	req := new(entity.RefreshTokenRequest)

	if err := c.Bind(req); err != nil {
		return nil, fmt.Errorf("%w: %s", entity.ErrorUnauthorized, err)
	}

	return req, nil
}

func parseRefreshTokenError(err error) (code int, i any) {
	switch {
	case errors.Is(err, entity.ErrorUnauthorized):
		return http.StatusUnauthorized, Response{Message: err.Error()}
	default:
		return http.StatusInternalServerError, Response{Message: "Internal server error"}
	}
}

func parseRefreshTokenResponse(accessToken *entity.AccessToken) (code int, i any) {
	return http.StatusOK, Response{Data: accessToken, Message: "Success refresh token"}
}

func parseRevokeTokenRequest(c echo.Context) (*entity.RevokeTokenRequest, error) {
	req := new(entity.RevokeTokenRequest)

	if err := c.Bind(req); err != nil {
		return nil, fmt.Errorf("%w: %s", entity.ErrorUnauthorized, err)
	}

	return req, nil
}

func parseRevokeTokenError(err error) (code int, i any) {
	switch {
	case errors.Is(err, entity.ErrorUnauthorized):
		return http.StatusUnauthorized, Response{Message: err.Error()}
	default:
		return http.StatusInternalServerError, Response{Message: "Internal server error"}
	}
}

func parseRevokeTokenResponse() (code int, i any) {
	return http.StatusOK, Response{Message: "Success revoke token"}
}

func parseGetUserByTokenRequest(c echo.Context) (*entity.GetUserByTokenRequest, error) {
	req := new(entity.GetUserByTokenRequest)

	if err := c.Bind(req); err != nil {
		return nil, fmt.Errorf("%w: %s", entity.ErrorUnauthorized, err)
	}

	return req, nil
}

func parseGetUserByTokenError(err error) (code int, i any) {
	switch {
	case errors.Is(err, entity.ErrorUnauthorized):
		return http.StatusUnauthorized, Response{Message: err.Error()}
	default:
		return http.StatusInternalServerError, Response{Message: "Internal server error"}
	}
}

func parseGetUserByTokenResponse(user *entity.User) (code int, i any) {
	return http.StatusOK, Response{Data: user, Message: "Success retrieve user"}
}
