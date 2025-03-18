package http_handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mocha-bot/mochus/config"
	"github.com/mocha-bot/mochus/core/module"
	cookiey "github.com/mocha-bot/mochus/pkg/cookie"
)

type discordHandler struct {
	cfg            config.Config
	discordUsecase module.DiscordUsecase
}

type DiscordHandler interface {
	OauthCallback(c echo.Context) error
	RefreshToken(c echo.Context) error
	RevokeToken(c echo.Context) error
	GetUserByToken(c echo.Context) error
}

func NewDiscordHandler(cfg config.Config, discordUsecase module.DiscordUsecase) DiscordHandler {
	return &discordHandler{
		cfg:            cfg,
		discordUsecase: discordUsecase,
	}
}

func (d *discordHandler) OauthCallback(c echo.Context) error {
	ctx := c.Request().Context()

	// If the error occurs, the discord will still got redirected to the desired URL
	// with the error message in the query params and the URI fragment will also be appended
	// e.g. /redirect?error=error_message#access_token=token&token_type=type&expires_in=3600&refresh_token=refresh&scope=identify

	req, err := parseOauthCallbackRequest(c)
	if err != nil {
		code, resp := parseOauthCallbackError(err)
		return c.Redirect(http.StatusTemporaryRedirect, parseOauthCallbackRedirectError(req.RedirectURL, code, resp))
	}

	exchanged, err := d.discordUsecase.ExchangeCodeForToken(ctx, req)
	if err != nil {
		code, resp := parseOauthCallbackError(err)
		return c.Redirect(http.StatusTemporaryRedirect, parseOauthCallbackRedirectError(req.RedirectURL, code, resp))
	}

	isLocalhost := d.cfg.App.IsLocalhost()

	for _, cookie := range exchanged.ToHTTPCookies() {
		cookie.Secure = !isLocalhost
		cookie.Domain = d.cfg.Discord.RedirectDomain
		cookie.Path = "/"
		cookie.SameSite = http.SameSiteLaxMode

		c.SetCookie(cookie)
	}

	return c.Redirect(http.StatusFound, req.RedirectURL)
}

func (d *discordHandler) RefreshToken(c echo.Context) error {
	ctx := c.Request().Context()

	req, err := parseRefreshTokenRequest(c)
	if err != nil {
		return c.JSON(parseRefreshTokenError(err))
	}

	exchanged, err := d.discordUsecase.ExchangeRefreshForToken(ctx, req)
	if err != nil {
		return c.JSON(parseRefreshTokenError(err))
	}

	isLocalhost := d.cfg.App.IsLocalhost()

	for _, cookie := range exchanged.ToHTTPCookies() {
		cookie.Secure = !isLocalhost
		cookie.Domain = d.cfg.Discord.RedirectDomain
		cookie.Path = "/"
		cookie.SameSite = http.SameSiteLaxMode

		c.SetCookie(cookie)
	}

	return c.JSON(http.StatusOK, Response{Message: "Token refreshed"})
}

func (d *discordHandler) RevokeToken(c echo.Context) error {
	ctx := c.Request().Context()

	req, err := parseRevokeTokenRequest(c)
	if err != nil {
		return c.JSON(parseRevokeTokenError(err))
	}

	err = d.discordUsecase.RevokeToken(ctx, req)
	if err != nil {
		return c.JSON(parseRevokeTokenError(err))
	}

	// remove cookies
	keys := []string{
		cookiey.CookieAccessToken,
		cookiey.CookieRefreshToken,
		cookiey.CookieTokenType,
		cookiey.CookieScope,
	}

	for _, key := range keys {
		c.SetCookie(&http.Cookie{Name: key, MaxAge: -1})
	}

	return c.JSON(http.StatusOK, Response{Message: "Token revoked"})
}

func (d *discordHandler) GetUserByToken(c echo.Context) error {
	ctx := c.Request().Context()

	req, err := parseGetUserByTokenRequest(c)
	if err != nil {
		return c.JSON(parseGetUserByTokenError(err))
	}

	user, err := d.discordUsecase.GetUser(ctx, req.ConstructAuthorization())
	if err != nil {
		return c.JSON(parseGetUserByTokenError(err))
	}

	return c.JSON(parseGetUserByTokenResponse(user))
}
