package cgroupfs

import (
	"fmt"
	"os"
	"path/filepath"
)

const CgroupRoot = "/sys/fs/cgroup"
const CgroupUserRoot = CgroupRoot + "/user.slice"

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

func UserSlices() ([]string, error) {
	entries, err := os.ReadDir(CgroupUserRoot)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	return names, nil
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
