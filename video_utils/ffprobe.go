package video_utils

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func GetVideoResolution(path string) (*VideoResolution, error) {
	// https://superuser.com/a/841379
	resolutionBytes, err := exec.Command(
		"ffprobe",
		"-v", "error", "-select_streams", "v:0", "-show_entries",
		"stream=width,height", "-of", "csv=s=x:p=0",
		path,
	).Output()

	if err != nil {
		return nil, fmt.Errorf("file not found")
	}

	// Resulting output is something like this: 1080x720
	resolutionStr := strings.ReplaceAll(
		string(resolutionBytes),
		"\n",
		"",
	)

	result := strings.Split(resolutionStr, "x")
	if len(result) != 2 {
		return nil, fmt.Errorf("invalid resolution: %s", resolutionStr)
	}

	width, err := strconv.Atoi(result[0])
	if err != nil {
		panic(err)
	}

	height, err := strconv.Atoi(result[1])
	if err != nil {
		panic(err)
	}

	videoResolution := CreateVideoResolution("", width, height)

	return &videoResolution, nil
}
