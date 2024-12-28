package cgroupfs

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

const CgroupRoot = "/sys/fs/cgroup"

func readUint64(path string) (uint64, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	var value uint64
	_, err = fmt.Fscanf(f, "%d", &value)
	return value, err
}

type CgroupInfo struct {
	Pids              uint64
	MemoryCurrent     uint64
	MemorySwapCurrent uint64
}

func GetCgroupInfo(dir string) (CgroupInfo, error) {
	info := CgroupInfo{}
	var err error
	info.Pids, err = readUint64(filepath.Join(dir, "pids.current"))
	if err != nil {
		return info, err
	}
	info.MemoryCurrent, err = readUint64(filepath.Join(dir, "memory.current"))
	if err != nil {
		return info, err
	}
	info.MemorySwapCurrent, err = readUint64(filepath.Join(dir, "memory.swap.current"))
	if err != nil {
		return info, err
	}
	return info, nil
}

func getAllPids(dir string) ([]int, []error) {
	var pids []int
	var errs []error
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, []error{err}
	}
	hasDir := false
	for _, entry := range entries {
		if entry.IsDir() {
			hasDir = true
			subpids, suberrs := getAllPids(filepath.Join(dir, entry.Name()))
			pids = append(pids, subpids...)
			errs = append(errs, suberrs...)
		}
	}
	if !hasDir {
		f, err := os.Open(filepath.Join(dir, "cgroup.procs"))
		if err != nil {
			return nil, []error{err}
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			pid, err := strconv.ParseInt(scanner.Text(), 10, 64)
			if err != nil {
				errs = append(errs, err)
				continue
			}
			pids = append(pids, int(pid))
		}
	}
	return pids, errs
}

func GetAllPids(dir string) ([]int, error) {
	pids, errs := getAllPids(dir)
	return pids, errors.Join(errs...)
}
