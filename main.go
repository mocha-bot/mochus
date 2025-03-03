package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mocha-bot/mochus/config"
	"github.com/mocha-bot/mochus/core/module"
	http_handler "github.com/mocha-bot/mochus/handler/http"
	discord_repository "github.com/mocha-bot/mochus/repository/discord"
	zLog "github.com/rs/zerolog/log"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		zLog.Fatal().Err(err).Msg("error loading config")
	}

	//need to add several instances several things such as (db, redis, etc)

	discordRepository := discord_repository.NewDiscordRepository(cfg.Discord)
	discordUsecase := module.NewDiscordUsecase(discordRepository)
	discordHandler := http_handler.NewDiscordHandler(cfg.Discord, discordUsecase)

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))
	e.Logger.SetLevel(log.LstdFlags)

	e.GET("/callback", discordHandler.OauthCallback)
	e.GET("/ridirect", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Hello, World!")
	})

	if err := e.Start(cfg.App.GetAddress()); err != nil {
		zLog.Fatal().Err(err).Msg("error starting server")
	}
}
