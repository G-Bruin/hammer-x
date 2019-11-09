package utils

import (
	"github.com/schollz/progressbar/v2"
	"time"
)

func Progressbar(length int64) *progressbar.ProgressBar {
	bar := progressbar.NewOptions(
		int(length),
		progressbar.OptionSetBytes(int(length)),
		progressbar.OptionThrottle(10*time.Millisecond),
	)
	return bar
}
