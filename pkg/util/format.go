package util

import (
	"fmt"
)

var sizeUnits = []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB"}

func FormatSize(size uint64) string {
	sizeF := float64(size)
	unit := 0
	for sizeF >= 1024 && unit < len(sizeUnits)-1 {
		sizeF /= 1024
		unit++
	}
	return fmt.Sprintf("%.1f %s", sizeF, sizeUnits[unit])
}
