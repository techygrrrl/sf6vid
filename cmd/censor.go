package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/techygrrrl/sf6vid/video_utils"
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

	// We use this to calculate the percentage-based censor boxes
	controlVideoResolution := video_utils.CreateVideoResolution("control", 1920, 1080)

	inputVideoResolution, err := video_utils.GetVideoResolution(inputPath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("video resolution: %s", inputVideoResolution)

	censorBoxes := []video_utils.CensorBox{
		video_utils.FixedSizeCensorBox{
			Name:   "Title",
			Width:  250,
			Height: 50,
			X:      300,
			Y:      8,
		}.ToCensorBox(controlVideoResolution),
		video_utils.FixedSizeCensorBox{
			Name:   "Rank and Club",
			Width:  190,
			Height: 115,
			X:      16,
			Y:      105,
		}.ToCensorBox(controlVideoResolution),
		video_utils.FixedSizeCensorBox{
			Name:   "Username",
			Width:  345,
			Height: 40,
			X:      205,
			Y:      106,
		}.ToCensorBox(controlVideoResolution),
	}

	for _, box := range censorBoxes {
		fmt.Println(box.PrettyJson())
	}
	fmt.Println("-------------------------")
	//filterComplex := ""

	if doP1 {
		fmt.Println("should censor player 1: ")
		for _, box := range censorBoxes {
			output, err := box.CropFilterOutput(*inputVideoResolution, video_utils.Player1)

			if err != nil {
				log.Fatal(err)
			}

			// TODO: remove - this is just for debugging
			fmt.Printf("Box: %s \nFilter output: %s \n\n", box.Name, output)
		}
		fmt.Println("-------------------------")
	}

	if doP2 {
		fmt.Println("should censor player 2: ")

		for _, box := range censorBoxes {
			output, err := box.CropFilterOutput(*inputVideoResolution, video_utils.Player2)

			if err != nil {
				log.Fatal(err)
			}

			// TODO: remove - this is just for debugging
			fmt.Printf("Box: %s \nFilter output: %s \n\n", box.Name, output)
		}
		fmt.Println("-------------------------")
	}
}
