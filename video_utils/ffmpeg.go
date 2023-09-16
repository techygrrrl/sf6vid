package video_utils

import (
	"fmt"
	"strings"
)

type ChainLink struct {
	CensorBox   CensorBox
	BlurSetting BlurSetting
}

func CreateChainLink(censorBox CensorBox, blurSetting BlurSetting) ChainLink {
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
		"[0:v]%s,%s;[base][tmp]%s[base]",
		cropFilterOutput,
		blurFilterOutput,
		//currentIndex,
		//currentIndex,
		overlayOutput,
	)

	return output, nil
}
