package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var flagUseDebug bool = false
var flagOpen bool = false

var rootCmd = &cobra.Command{
	Use:     "sf6vid",
	Version: "0.4.2",
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
	rootCmd.PersistentFlags().BoolVar(&flagOpen, "open", false, "Open the file after running this command")
}
