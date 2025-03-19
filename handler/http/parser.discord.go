package http_handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/mocha-bot/mochus/core/entity"
	cookiey "github.com/mocha-bot/mochus/pkg/cookiey"
	"github.com/mocha-bot/mochus/pkg/echoy"
	zLog "github.com/rs/zerolog/log"
)

func parseOauthCallbackRequest(c echo.Context) (*entity.OauthCallbackRequest, error) {
	req := new(entity.OauthCallbackRequest)

	if err := c.Bind(req); err != nil {
		return nil, fmt.Errorf("%w: %s", entity.ErrorBind, err)
	}

	parsedURL, err := url.ParseRequestURI(c.Request().RequestURI)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", entity.ErrorBind, err)
	}

	redirectURL := parsedURL.Query().Get(RedirectURLKey)

	// Construct the final URL for the request URL
	// Discord known this as a redirect_uri to verify the request
	finalURL := url.URL{
		Scheme: echoy.GetScheme(c),
		Host:   c.Request().Host,
		Path:   c.Request().URL.Path,
	}

	// The redirect URL is optional for the client redirection
	// If it's not provided, the fallback host will be used
	if redirectURL == "" {
		req.RedirectURL = c.Request().Header.Get("X-Fallback-Host")
	} else {
		finalURL.RawQuery = url.Values{RedirectURLKey: {req.RedirectURL}}.Encode()
	}

	req.RequestURL, err = url.QueryUnescape(finalURL.String())
	if err != nil {
		return nil, fmt.Errorf("%w: %s", entity.ErrorBind, err)
	}

	return req, nil
}

func parseOauthCallbackError(err error) (code int, i any) {
	switch {
	case errors.Is(err, entity.ErrorBind):
		return http.StatusBadRequest, Response{Message: err.Error()}
	case errors.Is(err, entity.ErrorUnauthorized):
		return http.StatusUnauthorized, Response{Message: err.Error()}
	case errors.Is(err, entity.ErrorBadRequest):
		return http.StatusBadRequest, Response{Message: err.Error()}
	default:
		zLog.Error().Err(err).Msg("Internal server error")
		return http.StatusInternalServerError, Response{Message: "Internal server error"}
	}
}

func parseOauthCallbackRedirectError(redirectURL string, code int, i any) (newRedirectURL string) {
	redirectURLParsed, _ := url.Parse(redirectURL)

	if i == nil {
		return redirectURLParsed.String()
	}

	redirectURLParsed.RawQuery = url.Values{"error": {i.(Response).Message}}.Encode()
	return redirectURLParsed.String()
}

func parseRefreshTokenRequest(c echo.Context) (*entity.RefreshTokenRequest, error) {
	req := new(entity.RefreshTokenRequest)

	refreshTokenCookie, err := c.Cookie(cookiey.CookieRefreshToken)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", entity.ErrorUnauthorized, err)
	}

	req.RefreshToken = refreshTokenCookie.Value

	return req, nil
}

func parseRefreshTokenError(err error) (code int, i any) {
	switch {
	case errors.Is(err, entity.ErrorUnauthorized):
		return http.StatusUnauthorized, Response{Message: err.Error()}
	case errors.Is(err, entity.ErrorBadRequest):
		return http.StatusBadRequest, Response{Message: err.Error()}
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

	if refreshTokenCookie != nil {
		req.RefreshToken = refreshTokenCookie.Value
	}

	if accessTokenCookie != nil {
		req.AccessToken = accessTokenCookie.Value
	}

	return req, nil
}

func parseRevokeTokenError(err error) (code int, i any) {
	switch {
	case errors.Is(err, entity.ErrorUnauthorized):
		return http.StatusUnauthorized, Response{Message: err.Error()}
	case errors.Is(err, entity.ErrorBadRequest):
		return http.StatusBadRequest, Response{Message: err.Error()}
	default:
		zLog.Error().Err(err).Msg("Internal server error")
		return http.StatusInternalServerError, Response{Message: "Internal server error"}
	}
}

func parseGetUserByTokenRequest(c echo.Context) (*entity.GetUserByTokenRequest, error) {
	req := new(entity.GetUserByTokenRequest)

	accessTokenCookie, err := c.Cookie(cookiey.CookieAccessToken)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", entity.ErrorUnauthorized, err)
	}

	req.AccessToken = accessTokenCookie.Value

	tokenTypeCookie, err := c.Cookie(cookiey.CookieTokenType)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", entity.ErrorUnauthorized, err)
	}

	req.TokenType = tokenTypeCookie.Value

	return req, nil
}

func parseGetUserByTokenError(err error) (code int, i any) {
	switch {
	case errors.Is(err, entity.ErrorUnauthorized):
		return http.StatusUnauthorized, Response{Message: err.Error()}
	case errors.Is(err, entity.ErrorBadRequest):
		return http.StatusBadRequest, Response{Message: err.Error()}
	default:
		zLog.Error().Err(err).Msg("Internal server error")
		return http.StatusInternalServerError, Response{Message: "Internal server error"}
	}
}

func parseGetUserByTokenResponse(user *entity.User) (code int, i any) {
	return http.StatusOK, Response{Data: user, Message: "Success retrieve user"}
}
