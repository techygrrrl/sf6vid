package string_utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppendStringToFileName_worksForMp4(t *testing.T) {
	filename := "out.mp4"
	stringToAppend := "640x360"

	result := AppendStringToFileName(filename, stringToAppend)

	assert.Equal(t, "out_640x360.mp4", result)
}

func TestAppendStringToFileName_worksForMov(t *testing.T) {
	filename := "out.mov"
	stringToAppend := "640x360"

	result := AppendStringToFileName(filename, stringToAppend)

	assert.Equal(t, "out_640x360.mov", result)
}
