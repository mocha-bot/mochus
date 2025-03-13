package main

import (
	"os"

	"github.com/mocha-bot/mochus/cmd"
	zLog "github.com/rs/zerolog/log"
)

func main() {
	err := cmd.RootCmd.Execute()
	if err != nil {
		zLog.Error().Err(err).Msg("failed to execute root command")
		os.Exit(1)
	}
}
