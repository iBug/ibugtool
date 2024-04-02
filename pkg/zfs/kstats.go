package zfs

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const KStatBase = "/proc/spl/kstat/zfs"

type ARCStats struct {
	Hits, Misses     uint64
	L2Hits, L2Misses uint64
}

func GetARCStats() (ARCStats, error) {
	f, err := os.Open(filepath.Join(KStatBase, "arcstats"))
	if err != nil {
		return ARCStats{}, err
	}
	defer f.Close()
	var s ARCStats
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) != 3 {
			continue
		}
		switch fields[0] {
		case "hits":
			s.Hits, err = strconv.ParseUint(fields[2], 10, 64)
		case "misses":
			s.Misses, err = strconv.ParseUint(fields[2], 10, 64)
		case "l2_hits":
			s.L2Hits, err = strconv.ParseUint(fields[2], 10, 64)
		case "l2_misses":
			s.L2Misses, err = strconv.ParseUint(fields[2], 10, 64)
		}
		if err != nil {
			return ARCStats{}, err
		}
	}
	return s, nil
}
