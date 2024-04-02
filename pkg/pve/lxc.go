package pve

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	ErrProcessNotLXC = errors.New("process does not belong to a container")
)

func FindCTByPid(pid int) (int, error) {
	procfile := fmt.Sprintf("/proc/%d/cgroup", pid)
	cmdline, err := os.ReadFile(procfile)
	if err != nil {
		return 0, err
	}
	parts := strings.Split(string(cmdline), "/")
	if len(parts) >= 3 && strings.HasSuffix(parts[0], "::") && parts[1] == "lxc" {
		return strconv.Atoi(parts[2])
	}
	return 0, ErrProcessNotLXC
}
