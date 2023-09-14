package video_math

import (
	"fmt"
	"math"
)

// region Players

type PlayerSide int8

const (
	Player1 PlayerSide = iota
	Player2
)

// endregion Players

// region Censor Boxes

type CensorBox struct {
	name             string
	widthPercentage  float64
	heightPercentage float64
	xPercentage      float64
	yPercentage      float64
}

func (c CensorBox) CropFilterOutput(v VideoResolution, side PlayerSide) (error, string) {
	cropWidth := int(math.Ceil(float64(v.width) * c.widthPercentage))
	cropHeight := int(math.Ceil(float64(v.height) * c.heightPercentage))
	cropY := int(math.Ceil(float64(v.height) * c.yPercentage))

	var cropX int = -1

	player1cropX := int(
		math.Ceil(
			float64(v.width) * c.xPercentage,
		),
	)
	// Player 2 side is a mirror of player 1, so we offset the X position accordingly
	var player2cropX int = int(
		math.Abs( // the resulting value is negative, so we use this to make it positive
			float64(
				player1cropX - (v.width - cropWidth),
			),
		),
	)

	if side == Player1 {
		cropX = player1cropX
	}

	if side == Player2 {
		cropX = player2cropX
	}

	if cropX == -1 {
		return fmt.Errorf("invalid player side %d", side), ""
	}

	return nil, fmt.Sprintf("crop=%d:%d:%d:%d", cropWidth, cropHeight, cropX, cropY)
}

// endregion Censor Boxes

// region Video

type VideoResolution struct {
	name   string
	width  int
	height int
}

func CreateVideoResolution(name string, width int, height int) VideoResolution {
	return VideoResolution{name, width, height}
}

func (v VideoResolution) Name() string {
	return v.name
}

func (v VideoResolution) Width() int {
	return v.width
}

func (v VideoResolution) Height() int {
	return v.height
}

// endregion Video

// region Blur settings

type BlurSetting int

func CreateBlurSetting(value int) BlurSetting {
	return BlurSetting(value)
}

func (b BlurSetting) FilterOutput() string {
	return fmt.Sprintf("avgblur=%d", b)
}

// endregion Blur settings
