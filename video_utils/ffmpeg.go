package video_utils

import (
	"fmt"
	"strings"
	"time"
)

type ChainLink struct {
	CensorBox   CensorBox
	BlurSetting BlurSettings
}

func CreateChainLink(censorBox CensorBox, blurSetting BlurSettings) ChainLink {
	return ChainLink{censorBox, blurSetting}
}

type ChainAssembler struct {
	ChainLinks []ChainLink
}

func CreateChainAssembler(chainLinks []ChainLink) ChainAssembler {
	return ChainAssembler{chainLinks}
}

func (c ChainAssembler) AssembleChain(v VideoResolution, side PlayerSide) (string, error) {
	output := make([]string, len(c.ChainLinks)+2)
	output[0] = "[0:v]copy[base]"

	for i, chainLink := range c.ChainLinks {
		currentIndex := i + 1

		chainLinkOutput, err := chainLink.AssembleChainLink(currentIndex, v, side)
		if err != nil {
			return "", err
		}

		output[currentIndex] = chainLinkOutput
	}

	output[len(c.ChainLinks)+1] = "[0:a]acopy"

	stringOutput := strings.Join(output, ";")

	return stringOutput, nil
}

func (c ChainLink) AssembleChainLink(currentIndex int, v VideoResolution, side PlayerSide) (string, error) {
	cropFilterOutput, err := c.CensorBox.CropFilterOutput(v, side)
	if err != nil {
		return "", err
	}

	overlayOutput, err := c.CensorBox.OverlayOutput(v, side)
	if err != nil {
		return "", err
	}

	blurFilterOutput := c.BlurSetting.FilterOutput()

	output := fmt.Sprintf(
		"[0:v]%s,%s[blur%d];[base][blur%d]%s[base]",
		cropFilterOutput,
		blurFilterOutput,
		currentIndex,
		currentIndex,
		overlayOutput,
	)

	return output, nil
}

func FormatDurationForFFmpeg(duration time.Duration) string {
	d := duration.Seconds()

	hour := int(d / 3600)
	minute := int(d/60) % 60
	second := int(d) % 60

	return fmt.Sprintf("%02d:%02d:%02d", hour, minute, second)
}

func FormattedDurationArgsForFFmpeg(start time.Duration, end time.Duration) []string {
	// if there is neither a start or end time, we assume the user did not provide any trimming configuration
	if start == 0 && end == 0 {
		return []string{}
	}

	var args []string
	if start != 0 {
		// configure start time
		args = []string{"-ss", FormatDurationForFFmpeg(start)}
	}

	if end != 0 {
		// configure end time
		args = append(args, "-to", FormatDurationForFFmpeg(end))
	}

	return args
}
