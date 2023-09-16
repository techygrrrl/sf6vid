package cmd

import (
	"fmt"
	"log"
	"os/exec"

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
	var playerSide video_utils.PlayerSide = -1
	if doP1 == doP2 {
		log.Fatalf("must specify only one of --p1 or --p2")
	}
	if doP1 {
		playerSide = video_utils.Player1
	}
	if doP2 {
		playerSide = video_utils.Player2
	}

	// We use this to calculate the percentage-based censor boxes
	controlVideoResolution := video_utils.CreateVideoResolution("control", 1920, 1080)

	inputVideoResolution, err := video_utils.GetVideoResolution(inputPath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("video resolution: %v", inputVideoResolution)

	// TODO: build this slice based on the p1 and p2 flags
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

	// TODO: dynamic blur setting
	blurSetting := video_utils.BlurSetting(4)

	chainLinks := make([]video_utils.ChainLink, len(censorBoxes))
	for i, box := range censorBoxes {
		chainLink := video_utils.CreateChainLink(box, blurSetting)
		chainLinks[i] = chainLink
	}

	chainAssembler := video_utils.CreateChainAssembler(chainLinks)

	filterComplexChain, err := chainAssembler.AssembleChain(*inputVideoResolution, playerSide)

	censorCommandOutput, err := exec.Command(
		"ffmpeg",
		"-i", inputPath,
		"-filter_complex", filterComplexChain,
		"-map", "[base]",
		//"-map \"[base]\"",
		"-y",
		outputPath,
	).Output()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("succes???? \n\n\n%s", string(censorCommandOutput))

	//ffmpeg \
	//-i samples/sample-720x480.mp4 \
	//-filter_complex "[0:v]crop=94:23:113:4,boxblur=4[blur1];[0:v][blur1]overlay=113:4[blurred1];[0:v]crop=72:52:6:47,boxblur=4[blur2];[blurred1][blur2]overlay=6:47[blurred2];[0:v]crop=130:18:77:48,boxblur=4[blur3];[blurred2][blur3]overlay=77:48" \
	//-y \
	//samples/blurred.mp4

	//currentIdx := 1
	//filterComplex := ""

	// Censor player 1
	//if doP1 {
	//	fmt.Println("should censor player 1: ")
	//	for _, box := range censorBoxes {
	//		output, err := box.CropFilterOutput(*inputVideoResolution, video_utils.Player1)
	//
	//		if err != nil {
	//			log.Fatal(err)
	//		}
	//
	//		// TODO: remove - this is just for debugging
	//		fmt.Printf("Box: %s \nFilter output: %s \n\n", box.Name, output)
	//	}
	//	fmt.Println("-------------------------")
	//}
	//
	//// Censor player 2
	//if doP2 {
	//	fmt.Println("should censor player 2: ")
	//
	//	for _, box := range censorBoxes {
	//		output, err := box.CropFilterOutput(*inputVideoResolution, video_utils.Player2)
	//
	//		if err != nil {
	//			log.Fatal(err)
	//		}
	//
	//		// TODO: remove - this is just for debugging
	//		fmt.Printf("Box: %s \nFilter output: %s \n\n", box.Name, output)
	//	}
	//	fmt.Println("-------------------------")
	//}
}
