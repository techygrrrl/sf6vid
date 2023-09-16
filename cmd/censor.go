package cmd

import (
	"fmt"
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
var blurSetting int

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
	censorCmd.Flags().IntVarP(&blurSetting, "blur", "b", 6, "Custom blur value")

	err := censorCmd.MarkFlagRequired("input")
	if err != nil {
		panic(err)
	}
	err = censorCmd.MarkFlagRequired("output")
	if err != nil {
		panic(err)
	}

	rootCmd.AddCommand(censorCmd)
}

func runCensorCmd(cmd *cobra.Command, args []string) {
	// Validation
	var playerSide video_utils.PlayerSide = -1
	if doP1 == doP2 {
		panic("must specify only one of --p1 or --p2")
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
		panic(err)
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

	chainLinks := make([]video_utils.ChainLink, len(censorBoxes))
	for i, box := range censorBoxes {
		chainLink := video_utils.CreateChainLink(box, video_utils.BlurSetting(blurSetting))
		chainLinks[i] = chainLink
	}

	chainAssembler := video_utils.CreateChainAssembler(chainLinks)

	filterComplexChain, err := chainAssembler.AssembleChain(*inputVideoResolution, playerSide)

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
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
		fmt.Printf("💥 could not process the video. try lowering the blur value from %d\n", blurSetting)
		os.Exit(1)
	}

	fullFilePath := fmt.Sprintf("%s/%s", cwd, outputPath)
	fmt.Printf("✅ Censored video should be available at %s\n", fullFilePath)

	if openFile {
		err = file_utils.OpenFile(fullFilePath)
	}
}
