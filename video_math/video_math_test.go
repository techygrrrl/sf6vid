package video_math

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var bigVideo = CreateVideoResolution("Big Video", 1920, 1080)
var smallVideo = CreateVideoResolution("Small Video", 960, 540)

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
	subject1 := CreateVideoResolution("SD", 720, 480)

	assert.Equal(t, "SD", subject1.Name())
	assert.Equal(t, 720, subject1.Width())
	assert.Equal(t, 480, subject1.Height())

	subject2 := CreateVideoResolution("beeg", 1920, 1080)

	assert.Equal(t, "beeg", subject2.Name())
	assert.Equal(t, 1920, subject2.Width())
	assert.Equal(t, 1080, subject2.Height())
}

func TestBlurSettings_FilterOutput(t *testing.T) {
	assert.Equal(t, "avgblur=10", CreateBlurSetting(10).FilterOutput())
	assert.Equal(t, "avgblur=20", CreateBlurSetting(20).FilterOutput())
}

func TestHardcodedCensorBox_ToCensorBox(t *testing.T) {
	title := HardcodedCensorBox{
		name:   "Title Box",
		width:  250,
		height: 50,
		x:      300,
		y:      8,
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
		HardcodedCensorBox{
			name:   "Title",
			width:  250,
			height: 50,
			x:      300,
			y:      8,
		}.ToCensorBox(bigVideo),
		HardcodedCensorBox{
			name:   "Rank and Club",
			width:  190,
			height: 115,
			x:      16,
			y:      105,
		}.ToCensorBox(bigVideo),
		HardcodedCensorBox{
			name:   "Username",
			width:  345,
			height: 40,
			x:      205,
			y:      106,
		}.ToCensorBox(bigVideo),
	}

	for _, box := range censorBoxes {
		fmt.Println(box.PrettyJson())
	}
}
