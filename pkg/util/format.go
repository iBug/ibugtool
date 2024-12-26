package util

import (
	"fmt"
)

var sizeUnits = []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB"}

func FormatSize(size uint64) string {
	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	}

	sizeF := float64(size)
	unit := 0
	for sizeF >= 1024 && unit < len(sizeUnits)-1 {
		sizeF /= 1024
		unit++
	}
	return fmt.Sprintf("%.1f %s", sizeF, sizeUnits[unit])
}

func FormatSizeAligned(size uint64) string {
	if size < 1024 {
		return fmt.Sprintf("%d     B", size)
	}
	return FormatSize(size)
}
