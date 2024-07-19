package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/techygrrrl/sf6vid/file_utils"
	"github.com/techygrrrl/sf6vid/string_utils"
	"github.com/techygrrrl/sf6vid/video_utils"
)

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
	censorCmd.Flags().SortFlags = false

	// player config
	censorCmd.Flags().Bool("p1", false, "Censor player 1 side")
	censorCmd.Flags().Bool("p2", false, "Censor player 2 side")

	// files
	censorCmd.Flags().StringP("input", "i", "", "Path to input file")
	censorCmd.Flags().StringP("output", "o", "", "Path to output file")

	// blur config
	censorCmd.Flags().Int("blur", 6, "Custom blur value for when the box blur is used (requires --box-blur flag otherwise this value will be ignored)")
	censorCmd.Flags().Bool("box-blur", false, "Use the box blur filter instead of the new pixelize filter (pixelize requires ffmpeg 6+)")

	// trim config
	censorCmd.Flags().Duration("start", time.Duration(0), "Optional start time for trimming the video")
	censorCmd.Flags().Duration("end", time.Duration(0), "Optional end time for trimming the video")

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
	doP1, err := cmd.Flags().GetBool("p1")
	if err != nil {
		panic(err)
	}

	doP2, err := cmd.Flags().GetBool("p2")
	if err != nil {
		panic(err)
	}

	inputPath, err := cmd.Flags().GetString("input")
	if err != nil {
		panic(err)
	}

	outputPath, err := cmd.Flags().GetString("output")
	if err != nil {
		panic(err)
	}

	blurValue, err := cmd.Flags().GetInt("blur")
	if err != nil {
		panic(err)
	}

	shouldUseLegacyBlur, err := cmd.Flags().GetBool("box-blur")
	if err != nil {
		panic(err)
	}

	startTime, err := cmd.Flags().GetDuration("start")
	if err != nil {
		panic(err)
	}

	endTime, err := cmd.Flags().GetDuration("end")
	if err != nil {
		panic(err)
	}

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
	controlVideoResolution := video_utils.CreateVideoResolution(1920, 1080)

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
			Height: 120,
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
		chainLink := video_utils.CreateChainLink(box, video_utils.CreateBlurSetting(blurValue, !shouldUseLegacyBlur))
		chainLinks[i] = chainLink
	}

	chainAssembler := video_utils.CreateChainAssembler(chainLinks)

	filterComplexChain, err := chainAssembler.AssembleChain(*inputVideoResolution, playerSide)
	if err != nil {
		panic(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	commandArgs := []string{
		"-i", inputPath,
		"-filter_complex", filterComplexChain,
		"-map", "[base]",

		// Quality settings. See shrink
		"-c:v", "libx265", "-crf", "30",

		"-y",
	}

	durationArgs := video_utils.FormattedDurationArgsForFFmpeg(startTime, endTime)

	commandArgs = append(commandArgs, durationArgs...)

	// append the output path
	var censoredPlayerString string
	if doP1 {
		censoredPlayerString = "p1"
	} else {
		censoredPlayerString = "p2"
	}

	durationSuffix := fmt.Sprintf("censored_%s_%s-%s", censoredPlayerString, startTime.String(), endTime.String())
	outputPathWithCensoredSuffix := string_utils.AppendStringToFileName(outputPath, durationSuffix)
	commandArgs = append(commandArgs, outputPathWithCensoredSuffix)

	if flagUseDebug {
		fmt.Printf("⚙️  Executing command:\n\nffmpeg %s\n\n", strings.Join(commandArgs, " "))
	}
	_, err = exec.Command("ffmpeg", commandArgs...).Output()
	if err != nil {
		fmt.Printf("💥 could not process the video. try lowering the blur value from %d\n", blurValue)
		os.Exit(1)
	}

	fullFilePath := fmt.Sprintf("%s/%s", cwd, outputPathWithCensoredSuffix)
	fmt.Printf("✅ Censored video was output to: %s\n", fullFilePath)

	if flagOpen {
		err = file_utils.OpenFile(fullFilePath)
	}
}
