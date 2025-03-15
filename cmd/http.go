package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mocha-bot/mochus/config"
	"github.com/mocha-bot/mochus/core/module"
	http_handler "github.com/mocha-bot/mochus/handler/http"
	discord_repository "github.com/mocha-bot/mochus/repository/discord"
	zLog "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var HTTPCmd = &cobra.Command{
	Use:   "http",
	Short: "Start listen to http server",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return serveHTTP(cmd, args)
	},
}

func serveHTTP(cmd *cobra.Command, args []string) error {
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}

	e := echo.New()
	e.Logger.SetLevel(log.LstdFlags)
	// e.Use(http_middleware.CORS(), middleware.RequestID())

	discordRepository := discord_repository.NewDiscordRepository(cfg.Discord)
	discordUsecase := module.NewDiscordUsecase(discordRepository)
	discordHandler := http_handler.NewDiscordHandler(cfg, discordUsecase)

	e.GET("/auth/discord/callback", discordHandler.OauthCallback)
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello World")
	})

	// Start server
	go func() {
		zLog.Info().Msgf("Starting server on %s", cfg.App.GetAddress())

		var err error
		if cfg.App.IsProduction() {
			// handled by cloudflare
			err = e.Start(cfg.App.GetAddress())
		} else {
			err = e.StartTLS(cfg.App.GetAddress(), "api-local.mocha-bot.xyz.pem", "api-local.mocha-bot.xyz-key.pem")
		}

		if err != nil {
			zLog.Fatal().Err(err).Msg("error starting server")
		}
	}()

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return e.Shutdown(ctx)
}
