package video_utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormatDuration(t *testing.T) {
	input, err := time.ParseDuration("4m8s")
	assert.Nil(t, err)

	result := FormatDuration(input)

	assert.Equal(t, "00:04:08", result)
}
