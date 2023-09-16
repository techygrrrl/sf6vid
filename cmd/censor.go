package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/techygrrrl/sf6vid/file_utils"
	"github.com/techygrrrl/sf6vid/video_utils"
)

var doP1 bool
var doP2 bool
var openFile bool
var inputPath string
var outputPath string

var censorCmd = &cobra.Command{
	Use:   "censor",
	Short: "Censor the player information in a video",
	Long: `Censor either the player 1 or player 2 identifying information in the video.
If the output path already exists, it will be replaced.
`,
	Run: runCensorCmd,
}

func init() {
	// Command options
	censorCmd.Flags().BoolVar(&doP1, "p1", false, "Censor player 1 side")
	censorCmd.Flags().BoolVar(&doP2, "p2", false, "Censor player 2 side")
	censorCmd.Flags().BoolVar(&openFile, "open", false, "Open the file after running this command")
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

	// TODO: dynamic blur setting
	blurSetting := video_utils.BlurSetting(4)

	chainLinks := make([]video_utils.ChainLink, len(censorBoxes))
	for i, box := range censorBoxes {
		chainLink := video_utils.CreateChainLink(box, blurSetting)
		chainLinks[i] = chainLink
	}

	chainAssembler := video_utils.CreateChainAssembler(chainLinks)

	filterComplexChain, err := chainAssembler.AssembleChain(*inputVideoResolution, playerSide)

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	_, err = exec.Command(
		"ffmpeg",
		"-i", inputPath,
		"-filter_complex", filterComplexChain,
		"-map", "[base]",
		"-y",
		outputPath,
	).Output()

	if err != nil {
		log.Fatal(err)
	}

	fullFilePath := fmt.Sprintf("%s/%s", cwd, outputPath)
	fmt.Printf("✅ Censored video should be available at %s\n", fullFilePath)

	if openFile {
		err = file_utils.OpenFile(fullFilePath)
	}
}
