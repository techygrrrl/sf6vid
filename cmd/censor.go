package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var doP1 bool
var doP2 bool
var inputPath string
var outputPath string

var censorCmd = &cobra.Command{
	Use:   "censor",
	Short: "Censor the player information in a video",
	Run:   runCensorCmd,
}

func init() {
	// Command options
	censorCmd.Flags().BoolVar(&doP1, "p1", false, "Censor player 1 side")
	censorCmd.Flags().BoolVar(&doP2, "p2", false, "Censor player 2 side")
	censorCmd.Flags().StringVarP(&inputPath, "input", "i", "", "Path to input file")
	censorCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Path to output file")

	err := censorCmd.MarkFlagRequired("input")
	if err != nil {
		log.Fatal(err)
	}
	err = censorCmd.MarkFlagRequired("output")
	if err != nil {
		log.Fatal(err)
	}

	rootCmd.AddCommand(censorCmd)
}

func runCensorCmd(cmd *cobra.Command, args []string) {
	fmt.Println("censor called")

	// Validation
	if !doP2 && !doP1 {
		log.Fatalf("must specify flags either --p1 or --p2 (or both)")
	}

	if doP1 {
		fmt.Println("should censor player 1")
	}

	if doP2 {
		fmt.Println("should censor player 2")
	}
}
