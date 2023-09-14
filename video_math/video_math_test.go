package video_math

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var bigVideo = CreateVideoResolution("Big Video", 1920, 1080)
var smallVideo = CreateVideoResolution("Small Video", 960, 540)

func TestCensorBox_CropFilterOutput_p1(t *testing.T) {
	// title box
	titleCensorBox := CensorBox{
		name:             "Title Box",
		widthPercentage:  0.130208333333333,
		heightPercentage: 0.046296296296296,
		xPercentage:      0.15625,
		yPercentage:      0.007407407407407,
	}

	err, smallResult := titleCensorBox.CropFilterOutput(smallVideo, Player1)
	assert.Nil(t, err)
	assert.Equal(t, "crop=125:25:150:4", smallResult)

	err, bigResult := titleCensorBox.CropFilterOutput(bigVideo, Player1)
	assert.Nil(t, err)
	assert.Equal(t, "crop=250:50:300:8", bigResult)
}

func TestCensorBox_CropFilterOutput_p2(t *testing.T) {
	// title box
	titleCensorBox := CensorBox{
		name:             "Title Box",
		widthPercentage:  0.130208333333333,
		heightPercentage: 0.046296296296296,
		xPercentage:      0.15625,
		yPercentage:      0.007407407407407,
	}

	err, smallResult := titleCensorBox.CropFilterOutput(smallVideo, Player2)
	assert.Nil(t, err)
	assert.Equal(t, "crop=125:25:685:4", smallResult)

	err, bigResult := titleCensorBox.CropFilterOutput(bigVideo, Player2)
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
