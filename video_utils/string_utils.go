package video_utils

import (
	"fmt"
	"time"
)

func FormatDuration(duration time.Duration) string {
	d := duration.Seconds()

	hour := int(d / 3600)
	minute := int(d/60) % 60
	second := int(d) % 60

	return fmt.Sprintf("%02d:%02d:%02d", hour, minute, second)
}
