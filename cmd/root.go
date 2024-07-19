package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var flagUseDebug bool = false

var rootCmd = &cobra.Command{
	Use:     "sf6vid",
	Version: "0.3.0",
}

func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&flagUseDebug, "debug", false, "More verbose logging")
}
