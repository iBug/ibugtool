package proc

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

const MeminfoFile = "/proc/meminfo"

type Meminfo struct {
	MemTotal uint64
}

func GetMeminfo() (Meminfo, error) {
	var info Meminfo
	f, err := os.Open(MeminfoFile)
	if err != nil {
		return info, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		switch fields[0] {
		case "MemTotal:":
			info.MemTotal, err = strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				return info, err
			}
			info.MemTotal *= 1024
		}
	}
	return info, nil
}
