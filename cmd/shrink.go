package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/techygrrrl/sf6vid/string_utils"
	"github.com/techygrrrl/sf6vid/video_utils"
)

var shrinkCmd = &cobra.Command{
	Use:   "shrink",
	Short: "Shrinks the video frame size by the indicated percentage, and compresses the quality in other ways",
	Long:  `Allows you to specify a percentage by which the video frame will be shrunk.`,
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
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

		targetPercent, err := cmd.Flags().GetInt("size")
		if err != nil {
			panic(err)
		}
		if targetPercent <= 0 || targetPercent >= 101 {
			panic(fmt.Errorf("invalid target percent: %d - `size` must be between 1 and 101", targetPercent))
		}

		inputVideoResolution, err := video_utils.GetVideoResolution(inputPath)
		if err != nil {
			panic(err)
		}

		scaledResolution := inputVideoResolution.GetScaledResolution(targetPercent)
		outputPathWithScaledSuffix := string_utils.AppendStringToFileName(outputPath, scaledResolution.String())

		commandArgs := []string{
			"-i", inputPath,

			// Quality-related. Constant Rate Factor, which lowers the average bit rate, but retains better quality.
			// Vary the CRF between around 18 and 24 — the lower, the higher the bitrate.
			// H.265 may use a crf between 24 to 30
			// See: https://unix.stackexchange.com/a/38380

			// This one is slightly larger but plays in QuickTime
			// "-b", "800k", "-crf", "28",
			// This one is slightly smaller but does not play in QuickTime but plays fine in VLC and Discord embeds
			"-c:v", "libx265", "-crf", "30",
			"-s", scaledResolution.String(),
			"-y",
			outputPathWithScaledSuffix,
		}

		if flagUseDebug {
			fmt.Printf("⚙️  Executing command:\n\nffmpeg %s\n\n", strings.Join(commandArgs, " "))
		}
		_, err = exec.Command("ffmpeg", commandArgs...).Output()
		if err != nil {
			panic(err)
		}

		fullFilePath := fmt.Sprintf("%s/%s", cwd, outputPathWithScaledSuffix)
		fmt.Printf("✅ Shrunk video was output to: %s\n", fullFilePath)
	},
}

func init() {
	shrinkCmd.Flags().SortFlags = false

	shrinkCmd.Flags().StringP("input", "i", "", "Path to input file")
	shrinkCmd.Flags().StringP("output", "o", "", "Path to output file")
	shrinkCmd.Flags().IntP("size", "s", 100, "Desired output size of the video percentage, e.g. a video that is 1280x720 will be 640x360 if you specify --size 50")

	rootCmd.AddCommand(shrinkCmd)
}
