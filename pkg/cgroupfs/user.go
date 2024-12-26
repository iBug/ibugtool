package cgroupfs

import (
	"os"
)

const CgroupUserRoot = CgroupRoot + "/user.slice"

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
