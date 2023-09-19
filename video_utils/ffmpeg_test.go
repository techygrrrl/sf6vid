package video_utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestChainLink_Assemble(t *testing.T) {
	currentIndex := 2
	bigVideo := CreateVideoResolution(1920, 1080)
	side := Player1

	subject := ChainLink{
		CensorBox: FixedSizeCensorBox{
			Name:   "Title",
			Width:  250,
			Height: 50,
			X:      300,
			Y:      8,
		}.ToCensorBox(bigVideo),
		BlurSetting: CreateBlurSetting(4, false),
	}

	result, err := subject.AssembleChainLink(currentIndex, bigVideo, side)

	assert.Nil(t, err)
	assert.Equal(t, "[0:v]crop=251:50:300:8,boxblur=4[blur2];[base][blur2]overlay=300:8[base]", result)
}

func TestChainAssembler_AssembleChain_blur(t *testing.T) {
	blurSetting := CreateBlurSetting(4, false)
	bigVideo := CreateVideoResolution(1920, 1080)
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
	assert.Equal(t, "[0:v]copy[base];[0:v]crop=251:50:300:8,boxblur=4[blur1];[base][blur1]overlay=300:8[base];[0:v]crop=190:115:16:105,boxblur=4[blur2];[base][blur2]overlay=16:105[base];[0:v]crop=345:40:205:106,boxblur=4[blur3];[base][blur3]overlay=205:106[base];[0:a]acopy", result)
}

func TestChainAssembler_AssembleChain_pixelize(t *testing.T) {
	blurSetting := CreateBlurSetting(4, true)
	bigVideo := CreateVideoResolution(1920, 1080)
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
	assert.Equal(t, "[0:v]copy[base];[0:v]crop=251:50:300:8,pixelize=16:16:avg[blur1];[base][blur1]overlay=300:8[base];[0:v]crop=190:115:16:105,pixelize=16:16:avg[blur2];[base][blur2]overlay=16:105[base];[0:v]crop=345:40:205:106,pixelize=16:16:avg[blur3];[base][blur3]overlay=205:106[base];[0:a]acopy", result)
}

func TestFormatDuration(t *testing.T) {
	input, err := time.ParseDuration("4m8s")
	assert.Nil(t, err)

	result := FormatDurationForFFmpeg(input)

	assert.Equal(t, "00:04:08", result)
}

func TestFormattedDurationArgsForFFmpeg(t *testing.T) {
	// with both start and end
	startInput1, err := time.ParseDuration("0m1s")
	assert.Nil(t, err)

	endInput1, err := time.ParseDuration("0m7s")
	assert.Nil(t, err)

	result1 := FormattedDurationArgsForFFmpeg(startInput1, endInput1)

	assert.Equal(t, []string{"-ss", "00:00:01", "-to", "00:00:07"}, result1)

	// with only start
	startInput2, err := time.ParseDuration("1m15s")
	assert.Nil(t, err)

	endInput2 := time.Duration(0)
	assert.Nil(t, err)

	result2 := FormattedDurationArgsForFFmpeg(startInput2, endInput2)

	assert.Equal(t, []string{"-ss", "00:01:15"}, result2)

	// with only end
	startInput3 := time.Duration(0)
	assert.Nil(t, err)

	endInput3, err := time.ParseDuration("3m30s")
	assert.Nil(t, err)

	result3 := FormattedDurationArgsForFFmpeg(startInput3, endInput3)
	assert.Equal(t, []string{"-to", "00:03:30"}, result3)
}
