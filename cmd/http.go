package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mocha-bot/mochus/config"
	"github.com/mocha-bot/mochus/core/module"
	http_handler "github.com/mocha-bot/mochus/handler/http"
	http_middleware "github.com/mocha-bot/mochus/handler/http/middleware"
	infrastructure_logger "github.com/mocha-bot/mochus/infrastructure/logger"
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

	zLog.Logger = infrastructure_logger.NewLogger(
		infrastructure_logger.WithConsole(cfg.Logger.ConsoleLogEnabled),
		infrastructure_logger.WithFile(cfg.Logger.FileLogEnabled),
		infrastructure_logger.WithLoggerConfig(&cfg.Logger),
	)

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Logger.SetOutput(zLog.Logger)

	CORS := http_middleware.CORS(
		http_middleware.WithAllowOrigins(cfg.App.CORSAllowOrigins),
		http_middleware.WithAllowMethods(cfg.App.CORSAllowMethods),
		http_middleware.WithAllowCredentials(cfg.App.CORSAllowCredentials),
	)

	e.Use(
		middleware.Recover(),
		middleware.RequestID(),
		middleware.Secure(),
		middleware.RequestLoggerWithConfig(http_middleware.RequestLoggerWithZerolog()),
		CORS,
		http_middleware.FallbackRedirect(cfg.App.FallbackRedirect),
	)

	discordRepository := discord_repository.NewDiscordRepository(cfg.Discord)
	discordUsecase := module.NewDiscordUsecase(discordRepository)
	discordHandler := http_handler.NewDiscordHandler(cfg, discordUsecase)

	apiV1 := e.Group("/api/v1")

	authRoute := apiV1.Group("/auth/discord")
	authRoute.GET("/callback", discordHandler.OauthCallback)
	authRoute.POST("/refresh", discordHandler.RefreshToken)
	authRoute.POST("/revoke", discordHandler.RevokeToken)
	authRoute.GET("/user", discordHandler.GetUserByToken)

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "mochus is healthy")
	})

	// Start server
	go func() {
		var err error

		switch {
		case cfg.App.IsTLS():
			if withCert := cfg.App.TLSCertFile != "" && cfg.App.TLSKeyFile != ""; withCert {
				zLog.Info().Msgf("Starting server on %s with TLS", cfg.App.GetAddress())
				err = e.StartTLS(cfg.App.GetAddress(), cfg.App.TLSCertFile, cfg.App.TLSKeyFile)
			} else {
				zLog.Info().Msgf("Starting server on %s with auto TLS", cfg.App.GetAddress())
				err = e.StartAutoTLS(cfg.App.GetAddress())
			}
		default:
			zLog.Info().Msgf("Starting server on %s", cfg.App.GetAddress())
			err = e.Start(cfg.App.GetAddress())
		}

		if err != nil && err != http.ErrServerClosed {
			zLog.Fatal().Err(err).Msg("error starting server")
		}
	}()

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	zLog.Info().Msg("Shutting down server")

	return e.Shutdown(ctx)
}
