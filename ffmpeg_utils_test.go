package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var bigVideo = videoResolution{
	name:   "Big Video",
	width:  1920,
	height: 1080,
}

var smallVideo = videoResolution{
	name:   "Small",
	width:  960,
	height: 540,
}

func TestCensorBox_CropFilterOutput(t *testing.T) {
	titleCensorBox := censorBox{
		name:             "Title Box",
		widthPercentage:  0.130208333333333,
		heightPercentage: 0.046296296296296,
		xPercentage:      0.15625,
		yPercentage:      0.007407407407407,
	}

	assert.Equal(t, "crop=125:25:150:4", titleCensorBox.CropFilterOutput(smallVideo))
	assert.Equal(t, "crop=250:50:300:8", titleCensorBox.CropFilterOutput(bigVideo))
}

func TestVideoResolution(t *testing.T) {
	subject1 := videoResolution{
		name:   "SD",
		width:  720,
		height: 480,
	}

	assert.Equal(t, "SD", subject1.Name())
	assert.Equal(t, 720, subject1.Width())
	assert.Equal(t, 480, subject1.Height())

	subject2 := videoResolution{
		name:   "beeg",
		width:  1920,
		height: 1080,
	}

	assert.Equal(t, "beeg", subject2.Name())
	assert.Equal(t, 1920, subject2.Width())
	assert.Equal(t, 1080, subject2.Height())
}

func TestBlurSettings_FilterOutput(t *testing.T) {
	assert.Equal(t, "avgblur=10", blurSettings{10}.FilterOutput())
	assert.Equal(t, "avgblur=20", blurSettings{20}.FilterOutput())
}
