package lteinfo

import (
	"time"
)

func ecomode(mode string) (time.Duration, time.Duration) {
	switch mode {
	case "off":
		return 5 * time.Millisecond, 1 * time.Second
	case "low":
		return 20 * time.Millisecond, 1 * time.Second
	case "mid":
		return 40 * time.Millisecond, 1 * time.Second
	case "normal":
		return 50 * time.Millisecond, 1 * time.Second
	case "high":
		return 80 * time.Millisecond, 3 * time.Second
	case "max":
		return 100 * time.Millisecond, 3 * time.Second
	case "extreme":
		return 200 * time.Millisecond, 5 * time.Second
	}
	return 35 * time.Millisecond, 1 * time.Second
}
