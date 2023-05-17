package utils

import (
	"fmt"
	"math"
)

func FormatReadableTime(time int64, showDays bool) string {
	hourSeconds := time
	if showDays {
		hourSeconds = time % (60 * 60 * 24)
	}
	minuteSeconds := hourSeconds % (60 * 60)
	remainingSeconds := minuteSeconds % 60

	days := int(math.Floor(float64(time / (60 * 60 * 24))))
	hours := int(math.Floor(float64(hourSeconds / (60 * 60))))
	minutes := int(math.Floor(float64(minuteSeconds / 60)))
	seconds := int(math.Floor(float64(remainingSeconds)))

	if days > 0 && showDays {
		return fmt.Sprintf("%dd %02dh %02dm %02ds", days, hours, minutes, seconds)
	}
	if hours > 0 {
		return fmt.Sprintf("%2dh %02dm %02ds", hours, minutes, seconds)
	}
	if minutes > 0 {
		return fmt.Sprintf("%2dm %02ds", minutes, seconds)
	}

	return fmt.Sprintf("%ds", seconds)
}

// Val Returns the value of the pointer or a default
func Val[T any](v *T, fallback T) T {
	if v != nil {
		return *v
	}
	return fallback
}
