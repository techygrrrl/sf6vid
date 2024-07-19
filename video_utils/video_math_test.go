package video_utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var bigVideo = CreateVideoResolution(1920, 1080)
var smallVideo = CreateVideoResolution(960, 540)

func TestCensorBox_CropFilterOutput_p1(t *testing.T) {
	// title box
	titleCensorBox := CensorBox{
		Name:             "Title Box",
		WidthPercentage:  0.130208333333333,
		HeightPercentage: 0.046296296296296,
		XPercentage:      0.15625,
		YPercentage:      0.007407407407407,
	}

	smallResult, err := titleCensorBox.CropFilterOutput(smallVideo, Player1)
	assert.Nil(t, err)
	assert.Equal(t, "crop=125:25:150:4", smallResult)

	bigResult, err := titleCensorBox.CropFilterOutput(bigVideo, Player1)
	assert.Nil(t, err)
	assert.Equal(t, "crop=250:50:300:8", bigResult)
}

func TestCensorBox_CropFilterOutput_p2(t *testing.T) {
	// title box
	titleCensorBox := CensorBox{
		Name:             "Title Box",
		WidthPercentage:  0.130208333333333,
		HeightPercentage: 0.046296296296296,
		XPercentage:      0.15625,
		YPercentage:      0.007407407407407,
	}

	smallResult, err := titleCensorBox.CropFilterOutput(smallVideo, Player2)
	assert.Nil(t, err)
	assert.Equal(t, "crop=125:25:685:4", smallResult)

	bigResult, err := titleCensorBox.CropFilterOutput(bigVideo, Player2)
	assert.Nil(t, err)
	assert.Equal(t, "crop=250:50:1370:8", bigResult)
}

func TestVideoResolution(t *testing.T) {
	subject1 := CreateVideoResolution(720, 480)

	assert.Equal(t, 720, subject1.Width())
	assert.Equal(t, 480, subject1.Height())

	subject2 := CreateVideoResolution(1920, 1080)

	assert.Equal(t, 1920, subject2.Width())
	assert.Equal(t, 1080, subject2.Height())
}

func TestBlurSettings_FilterOutput(t *testing.T) {
	assert.Equal(t, "boxblur=10", CreateBlurSetting(10, false).FilterOutput())
	assert.Equal(t, "boxblur=20", CreateBlurSetting(20, false).FilterOutput())
}

func TestHardcodedCensorBox_ToCensorBox(t *testing.T) {
	title := FixedSizeCensorBox{
		Name:   "Title Box",
		Width:  250,
		Height: 50,
		X:      300,
		Y:      8,
	}
	expectedTitle := CensorBox{
		Name:             "Title Box",
		WidthPercentage:  0.130208333333333,
		HeightPercentage: 0.046296296296296,
		XPercentage:      0.15625,
		YPercentage:      0.007407407407407,
	}
	titleResult := title.ToCensorBox(bigVideo)

	fmt.Println(titleResult.PrettyJson())
	assert.InDelta(t, expectedTitle.WidthPercentage, titleResult.WidthPercentage, 0.00001)
	assert.InDelta(t, expectedTitle.HeightPercentage, titleResult.HeightPercentage, 0.00001)
	assert.InDelta(t, expectedTitle.XPercentage, titleResult.XPercentage, 0.00001)
	assert.InDelta(t, expectedTitle.YPercentage, titleResult.YPercentage, 0.00001)
}

func TestHardcodedCensorBox_ToCensorBox_moreBoxes(t *testing.T) {
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

	for _, box := range censorBoxes {
		fmt.Println(box.PrettyJson())
	}
}

func TestCensorBox_OverlayOutput(t *testing.T) {
	// title box
	titleCensorBox := CensorBox{
		Name:             "Title Box",
		WidthPercentage:  0.130208333333333,
		HeightPercentage: 0.046296296296296,
		XPercentage:      0.15625,
		YPercentage:      0.007407407407407,
	}

	smallResult, err := titleCensorBox.OverlayOutput(smallVideo, Player2)
	assert.Nil(t, err)
	assert.Equal(t, "overlay=685:4", smallResult)
}

func TestGetScaledResolution(t *testing.T) {
	subject := VideoResolution{
		width:  1280,
		height: 720,
	}

	output := subject.GetScaledResolution(50)
	// Correct values
	assert.Equal(t, 640, output.width)
	assert.Equal(t, 360, output.height)
	// Original struct unchanged
	assert.Equal(t, 1280, subject.width)
	assert.Equal(t, 720, subject.height)
	fmt.Printf("Scaled: %v , Original: %v ", output, subject)
}
