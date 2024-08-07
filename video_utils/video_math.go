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
	width  int
	height int
}

func CreateVideoResolution(width int, height int) VideoResolution {
	return VideoResolution{width, height}
}

func (v VideoResolution) Width() int {
	return v.width
}

func (v VideoResolution) Height() int {
	return v.height
}

func (v VideoResolution) String() string {
	return fmt.Sprintf("%dx%d", v.width, v.height)
}

func (v VideoResolution) GetScaledResolution(scalePercent int) VideoResolution {
	sizeDiff := float32(scalePercent) * 0.01
	targetWidth := float32(v.Width()) * sizeDiff
	targetHeight := float32(v.Height()) * sizeDiff

	return VideoResolution{
		width:  int(targetWidth),
		height: int(targetHeight),
	}
}

// endregion Video

// region Blur settings

type BlurSettings struct {
	BoxBlur  int
	Pixelize bool
}

func CreateBlurSetting(value int, pixelize bool) BlurSettings {
	return BlurSettings{value, pixelize}
}

func (b BlurSettings) FilterOutput() string {
	if b.Pixelize {
		return "pixelize=16:16:avg"
	}

	return fmt.Sprintf("boxblur=%d", b.BoxBlur)
}

// endregion Blur settings
