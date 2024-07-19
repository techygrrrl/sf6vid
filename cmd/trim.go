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

var trimCmd = &cobra.Command{
	Use:   "trim",
	Short: "Trim the video for the provided start and/or end times",
	Long: `You can provide one or both flags --start and --end.
If you omit --start, the original start time of the video will be used.
If you omit --end, the original end time of the video will be used.
At least one is required.
--start and --end use duration syntax, e.g. 5m30s for 5 minutes and 30 seconds
`,
	Run: runTrimCmd,
}

func init() {
	trimCmd.Flags().SortFlags = false

	// files
	trimCmd.Flags().StringP("input", "i", "", "Path to input file")
	trimCmd.Flags().StringP("output", "o", "", "Path to output file")

	// trim config
	trimCmd.Flags().Duration("start", time.Duration(0), "Start time for trimming the video")
	trimCmd.Flags().Duration("end", time.Duration(0), "End time for trimming the video")

	err := trimCmd.MarkFlagRequired("input")
	if err != nil {
		panic(err)
	}
	err = trimCmd.MarkFlagRequired("output")
	if err != nil {
		panic(err)
	}

	rootCmd.AddCommand(trimCmd)
}

func runTrimCmd(cmd *cobra.Command, args []string) {
	inputPath, err := cmd.Flags().GetString("input")
	if err != nil {
		panic(err)
	}

	outputPath, err := cmd.Flags().GetString("output")
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
	if startTime == 0 && endTime == 0 {
		fmt.Println("🤨 it looks like you didn't specify --start or --end so there will be no trimming. Exiting.")
		os.Exit(1)
	}

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	durationArgs := video_utils.FormattedDurationArgsForFFmpeg(startTime, endTime)
	commandArgs := []string{
		"-i", inputPath,

		// Quality settings. See shrink
		"-c:v", "libx265", "-crf", "30",

		"-y",
	}

	commandArgs = append(commandArgs, durationArgs...)

	durationSuffix := fmt.Sprintf("trimmed_%s-%s", startTime.String(), endTime.String())
	outputPathWithTrimmedSuffix := string_utils.AppendStringToFileName(outputPath, durationSuffix)
	commandArgs = append(commandArgs, outputPathWithTrimmedSuffix)

	if flagUseDebug {
		fmt.Printf("⚙️  Executing command:\n\nffmpeg %s\n\n", strings.Join(commandArgs, " "))
	}
	_, err = exec.Command("ffmpeg", commandArgs...).Output()
	if err != nil {
		fmt.Println("💥 could not trim the video")
		os.Exit(1)
	}

	fullFilePath := fmt.Sprintf("%s/%s", cwd, outputPathWithTrimmedSuffix)
	fmt.Printf("✅ Trimmed video was output to: %s\n", fullFilePath)

	if flagOpen {
		err = file_utils.OpenFile(fullFilePath)
	}
}
