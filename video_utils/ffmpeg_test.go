package video_utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChainLink_Assemble(t *testing.T) {
	currentIndex := 2
	bigVideo := CreateVideoResolution("big", 1920, 1080)
	side := Player1

	subject := ChainLink{
		CensorBox: FixedSizeCensorBox{
			Name:   "Title",
			Width:  250,
			Height: 50,
			X:      300,
			Y:      8,
		}.ToCensorBox(bigVideo),
		BlurSetting: BlurSetting(4),
	}

	result, err := subject.AssembleChainLink(currentIndex, bigVideo, side)

	assert.Nil(t, err)
	assert.Equal(t, "[0:v]crop=251:50:300:8,boxblur=4[blur2];[base][blur2]overlay=300:8[base]", result)
}

func TestChainAssembler_AssembleChain(t *testing.T) {
	blurSetting := BlurSetting(4)
	bigVideo := CreateVideoResolution("big", 1920, 1080)
	side := Player1
	censorBoxes := []CensorBox{
		FixedSizeCensorBox{
			Name:   "Title",
			Width:  250,
			Height: 50,
			X:      300,
			Y:      8,
		}.ToCensorBox(bigVideo),
		FixedSizeCensorBox{
			Name:   "Rank and Club",
			Width:  190,
			Height: 115,
			X:      16,
			Y:      105,
		}.ToCensorBox(bigVideo),
		FixedSizeCensorBox{
			Name:   "Username",
			Width:  345,
			Height: 40,
			X:      205,
			Y:      106,
		}.ToCensorBox(bigVideo),
	}
	chainLinks := make([]ChainLink, len(censorBoxes))
	for i, box := range censorBoxes {
		chainLink := CreateChainLink(box, blurSetting)
		chainLinks[i] = chainLink
	}

	subject := CreateChainAssembler(chainLinks)
	result, err := subject.AssembleChain(bigVideo, side)

	assert.Nil(t, err)
	assert.Equal(t, "[0:v]copy[base];[0:v]crop=251:50:300:8,boxblur=4[blur1];[base][blur1]overlay=300:8[base];[0:v]crop=190:115:16:105,boxblur=4[blur2];[base][blur2]overlay=16:105[base];[0:v]crop=345:40:205:106,boxblur=4[blur3];[base][blur3]overlay=205:106[base]", result)
}
