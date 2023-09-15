package video_utils

import (
	"encoding/json"
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
	Name             string
	WidthPercentage  float64
	HeightPercentage float64
	XPercentage      float64
	YPercentage      float64
}

func (c CensorBox) PrettyJson() (string, error) {
	asJson, err := json.MarshalIndent(c, "", "  ")

	if err != nil {
		return "", err
	}

	return string(asJson), nil
}

// HardcodedCensorBox To assist with calculating a CensorBox when provided a VideoResolution
type HardcodedCensorBox struct {
	Name   string
	Width  int
	Height int
	X      int
	Y      int
}

func (box HardcodedCensorBox) ToCensorBox(v VideoResolution) CensorBox {
	return CensorBox{
		Name:             box.Name,
		WidthPercentage:  float64(box.Width) / float64(v.width),
		HeightPercentage: float64(box.Height) / float64(v.height),
		XPercentage:      float64(box.X) / float64(v.width),
		YPercentage:      float64(box.Y) / float64(v.height),
	}
}

func (c CensorBox) CropFilterOutput(v VideoResolution, side PlayerSide) (string, error) {
	cropWidth := int(math.Ceil(float64(v.width) * c.WidthPercentage))
	cropHeight := int(math.Ceil(float64(v.height) * c.HeightPercentage))
	cropY := int(math.Ceil(float64(v.height) * c.YPercentage))

	var cropX int = -1

	player1cropX := int(
		math.Ceil(
			float64(v.width) * c.XPercentage,
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
		return "", fmt.Errorf("invalid player side %d", side)
	}

	return fmt.Sprintf("crop=%d:%d:%d:%d", cropWidth, cropHeight, cropX, cropY), nil
}

// endregion Censor Boxes

// region Video

type VideoResolution struct {
	name   string // todo: consider removing this field
	width  int
	height int
}

func CreateVideoResolution(name string, width int, height int) VideoResolution {
	return VideoResolution{name, width, height}
}

// Name todo: consider removing this field
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
