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
	//output := []string{
	//	"[0:v]copy[base]",
	//}

	output := make([]string, len(c.ChainLinks)+1)
	output[0] = "[0:v]copy[base]"

	//output := "[0:v]copy[base];"
	for i, chainLink := range c.ChainLinks {
		currentIndex := i + 1

		chainLinkOutput, err := chainLink.AssembleChainLink(currentIndex, v, side)
		if err != nil {
			return "", err
		}

		output[currentIndex] = chainLinkOutput
		//output = append(output, chainLinkOutput)
	}

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
	//CMD="$CMD[0:v]crop=94:23:0:0,boxblur=4[blur0];"
	//CMD="$CMD[base][blur0]overlay=113:4[base];"

	//output := fmt.Sprintf(
	//	"[0:v]%s,%s[blur%d];[%s][blur%d]%s[blurred%d]",
	//	cropFilterOutput,
	//	c.BlurSetting.FilterOutput(),
	//	currentIndex,
	//	c.OverlaySource,
	//	currentIndex,
	//	overlayOutput,
	//	currentIndex,
	//)

	return output, nil
}

//[0:v]crop=94:23:113:4,boxblur=0[blur1]; DONE
//[0:v][blur1]overlay=113:4[blurred1]; TODO
