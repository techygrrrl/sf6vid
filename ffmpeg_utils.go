package main

import (
	"fmt"
	"math"
)

// region Censor Boxes

type censorBox struct {
	name             string
	widthPercentage  float64
	heightPercentage float64
	xPercentage      float64
	yPercentage      float64
}

type PlayerSide int8

const (
	Player1 PlayerSide = iota
	Player2
)

func (c censorBox) CropFilterOutput(v videoResolution, side PlayerSide) (error, string) {
	if side == Player1 {
		cropWidth := int(math.Ceil(float64(v.width) * c.widthPercentage))
		cropHeight := int(math.Ceil(float64(v.height) * c.heightPercentage))
		cropX := int(math.Ceil(float64(v.width) * c.xPercentage))
		cropY := int(math.Ceil(float64(v.height) * c.yPercentage))

		return nil, fmt.Sprintf("crop=%d:%d:%d:%d", cropWidth, cropHeight, cropX, cropY)
	}

	if side == Player2 {
		// todo:
	}

	return fmt.Errorf("invalid player side %d", side), ""
}

// endregion Censor Boxes

// region Video

type videoResolution struct {
	name   string
	width  int
	height int
}

func (v videoResolution) Name() string {
	return v.name
}

func (v videoResolution) Width() int {
	return v.width
}

func (v videoResolution) Height() int {
	return v.height
}

// endregion Video

// region Blur settings

type blurSettings struct {
	value int
}

func (b blurSettings) FilterOutput() string {
	return fmt.Sprintf("avgblur=%d", b.value)
}

// endregion Blur settings
