package http_handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mocha-bot/mochus/config"
	"github.com/mocha-bot/mochus/core/module"
)

type discordHandler struct {
	cfg            config.DiscordConfig
	discordUsecase module.DiscordUsecase
}

type DiscordHandler interface {
	Oauth(c echo.Context) error
	OauthCallback(c echo.Context) error
}

func NewDiscordHandler(cfg config.DiscordConfig, discordUsecase module.DiscordUsecase) DiscordHandler {
	return &discordHandler{
		cfg:            cfg,
		discordUsecase: discordUsecase,
	}
}

func (d *discordHandler) Oauth(c echo.Context) error {
	return c.Redirect(http.StatusFound, d.cfg.GetOAuthURL())
}

func (d *discordHandler) OauthCallback(c echo.Context) error {
	code := c.QueryParam("code")
	ctx := c.Request().Context()
	if code == "" {
		return c.JSON(http.StatusBadRequest, "code is empty")
	}

	token, err := d.discordUsecase.ExchangeCodeForToken(ctx, code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Message: "Failed to exchange code for token",
		})
	}

	user, err := d.discordUsecase.GetUser(ctx, token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Message: "Failed to get user",
		})
	}

	return c.JSON(http.StatusOK, user)
}
