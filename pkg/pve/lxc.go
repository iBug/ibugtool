package pve

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	ErrProcessNotExist = errors.New("process does not exist")
	ErrInvalidCmdline  = errors.New("invalid command line for lxc-start")
)

func FindCTByPid(pid int) (int, error) {
	lastPid := 0
	for {
		procfile := fmt.Sprintf("/proc/%d/status", pid)
		f, err := os.Open(procfile)
		if errors.Is(err, os.ErrNotExist) {
			f.Close()
			return 0, ErrProcessNotExist
		}
		var name string
		var ppid int
		s := bufio.NewScanner(f)
		for s.Scan() {
			parts := strings.Fields(s.Text())
			if len(parts) < 2 {
				continue
			}
			switch parts[0] {
			case "Name:":
				name = parts[1]
			case "PPid:":
				ppid, err = strconv.Atoi(parts[1])
				if err != nil {
					return 0, err
				}
			}
		}
		f.Close()
		if name == "lxc-start" {
			lastPid = pid
		} else if ppid == 1 {
			break
		}
	}

	if lastPid == 0 {
		return 0, nil
	}
	procfile := fmt.Sprintf("/proc/%d/cmdline", pid)
	cmdline, err := os.ReadFile(procfile)
	if err != nil {
		return 0, err
	}
	parts := strings.Split(string(cmdline), "\000")
	for i, part := range parts {
		if part == "-n" {
			if i+1 >= len(parts) {
				return 0, ErrInvalidCmdline
			}
			return strconv.Atoi(parts[i+1])
		}
	}
	return 0, ErrInvalidCmdline
}
