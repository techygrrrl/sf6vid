package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/techygrrrl/sf6vid/video_math"
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

	// TODO: replace this with call to ffprobe
	inputVideoResolution := video_math.CreateVideoResolution("Video", 1920, 1080)

	censorBoxes := []video_math.CensorBox{
		video_math.HardcodedCensorBox{
			Name:   "Title",
			Width:  250,
			Height: 50,
			X:      300,
			Y:      8,
		}.ToCensorBox(inputVideoResolution),
		video_math.HardcodedCensorBox{
			Name:   "Rank and Club",
			Width:  190,
			Height: 115,
			X:      16,
			Y:      105,
		}.ToCensorBox(inputVideoResolution),
		video_math.HardcodedCensorBox{
			Name:   "Username",
			Width:  345,
			Height: 40,
			X:      205,
			Y:      106,
		}.ToCensorBox(inputVideoResolution),
	}
	// TODO: remove - this is just for debugging
	for _, box := range censorBoxes {
		fmt.Println(box.PrettyJson())
	}
	fmt.Println("-------------------------")

	if doP1 {
		fmt.Println("should censor player 1: ")
		for _, box := range censorBoxes {
			output, err := box.CropFilterOutput(inputVideoResolution, video_math.Player1)

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
			output, err := box.CropFilterOutput(inputVideoResolution, video_math.Player2)

			if err != nil {
				log.Fatal(err)
			}

			// TODO: remove - this is just for debugging
			fmt.Printf("Box: %s \nFilter output: %s \n\n", box.Name, output)
		}
		fmt.Println("-------------------------")
	}
}
