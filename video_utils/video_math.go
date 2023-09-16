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

// FixedSizeCensorBox To assist with calculating a CensorBox when provided a VideoResolution
type FixedSizeCensorBox struct {
	Name   string
	Width  int
	Height int
	X      int
	Y      int
}

func (box FixedSizeCensorBox) ToCensorBox(v VideoResolution) CensorBox {
	return CensorBox{
		Name:             box.Name,
		WidthPercentage:  float64(box.Width) / float64(v.width),
		HeightPercentage: float64(box.Height) / float64(v.height),
		XPercentage:      float64(box.X) / float64(v.width),
		YPercentage:      float64(box.Y) / float64(v.height),
	}
}

func (c CensorBox) GetYPositionForPlayerSide(v VideoResolution, side PlayerSide) int {
	return int(math.Ceil(float64(v.height) * c.YPercentage))
}

func (c CensorBox) GetXPositionForPlayerSide(v VideoResolution, side PlayerSide) (int, error) {
	cropWidth := int(math.Ceil(float64(v.width) * c.WidthPercentage))

	player1cropX := int(
		math.Ceil(
			float64(v.width) * c.XPercentage,
		),
	)
	// Player 2 side is a mirror of player 1, so we offset the X position accordingly
	player2cropX := int(
		math.Abs( // the resulting value is negative, so we use this to make it positive
			float64(
				player1cropX - (v.width - cropWidth),
			),
		),
	)

	if side == Player1 {
		return player1cropX, nil
	}

	if side == Player2 {
		return player2cropX, nil
	}

	return -1, fmt.Errorf("invalid player side: %d", side)
}

func (c CensorBox) GetWidthForPlayerSide(v VideoResolution, side PlayerSide) int {
	return int(math.Ceil(float64(v.width) * c.WidthPercentage))
}

func (c CensorBox) GetHeightForPlayerSide(v VideoResolution, side PlayerSide) int {
	return int(math.Ceil(float64(v.height) * c.HeightPercentage))
}

func (c CensorBox) CropFilterOutput(v VideoResolution, side PlayerSide) (string, error) {
	cropX, err := c.GetXPositionForPlayerSide(v, side)
	if err != nil {
		return "", err
	}

	cropWidth := c.GetWidthForPlayerSide(v, side)
	cropHeight := c.GetHeightForPlayerSide(v, side)
	cropY := c.GetYPositionForPlayerSide(v, side)

	cropFilterOutput := fmt.Sprintf("crop=%d:%d:%d:%d", cropWidth, cropHeight, cropX, cropY)

	return cropFilterOutput, nil
}

func (c CensorBox) OverlayOutput(v VideoResolution, side PlayerSide) (string, error) {
	cropX, err := c.GetXPositionForPlayerSide(v, side)
	if err != nil {
		return "", err
	}

	cropY := c.GetYPositionForPlayerSide(v, side)

	overlayOutput := fmt.Sprintf("overlay=%d:%d", cropX, cropY)

	return overlayOutput, nil
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
	// todo: pass the CensorBox here as an argument and get the Math.min value
	return BlurSetting(value)
}

func (b BlurSetting) FilterOutput() string {
	// todo: pass the CensorBox here as an argument and get the Math.min value
	return fmt.Sprintf("avgblur=%d", b)
	//return fmt.Sprintf("boxblur=%d", b)
}

// endregion Blur settings
