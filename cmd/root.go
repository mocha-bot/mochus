package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "mochus",
	Short: "Mocha Application Programming Interface (API) - User Service",
}

func init() {
	cobra.OnInitialize()

	// Initialize commands
	RootCmd.AddCommand(HTTPCmd)
}
