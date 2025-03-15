package http_handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mocha-bot/mochus/config"
	"github.com/mocha-bot/mochus/core/module"
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

	req, err := parseOauthCallbackRequest(c)
	if err != nil {
		return c.JSON(parseOauthCallbackError(err))
	}

	exchanged, err := d.discordUsecase.ExchangeCodeForToken(ctx, req.Code)
	if err != nil {
		return c.JSON(parseOauthCallbackError(err))
	}

	isLocalhost := d.cfg.App.IsLocalhost()

	for _, cookie := range exchanged.ToHTTPCookies() {
		cookie.Secure = !isLocalhost
		cookie.Domain = d.cfg.Discord.RedirectDomain
		cookie.Path = "/"
		cookie.SameSite = http.SameSiteLaxMode

		c.SetCookie(cookie)
	}

	return c.Redirect(http.StatusFound, d.cfg.App.Gateway)
}

func (d *discordHandler) RefreshToken(c echo.Context) error {
	ctx := c.Request().Context()

	req, err := parseRefreshTokenRequest(c)
	if err != nil {
		return c.JSON(parseRefreshTokenError(err))
	}

	exchanged, err := d.discordUsecase.ExchangeRefreshForToken(ctx, req.RefreshToken)
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

	return c.Redirect(http.StatusFound, d.cfg.App.Gateway)
}

func (d *discordHandler) RevokeToken(c echo.Context) error {
	ctx := c.Request().Context()

	req, err := parseRevokeTokenRequest(c)
	if err != nil {
		return c.JSON(parseRevokeTokenError(err))
	}

	err = d.discordUsecase.RevokeToken(ctx, req.Token)
	if err != nil {
		return c.JSON(parseRevokeTokenError(err))
	}

	return c.JSON(http.StatusOK, Response{Message: "Success revoke token"})
}

func (d *discordHandler) GetUserByToken(c echo.Context) error {
	ctx := c.Request().Context()

	req, err := parseGetUserByTokenRequest(c)
	if err != nil {
		return c.JSON(parseGetUserByTokenError(err))
	}

	user, err := d.discordUsecase.GetUser(ctx, req.Token())
	if err != nil {
		return c.JSON(parseGetUserByTokenError(err))
	}

	return c.JSON(parseGetUserByTokenResponse(user))
}
