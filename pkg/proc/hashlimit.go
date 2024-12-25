package proc

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const HashlimitBasePath = "/proc/net/ipt_hashlimit"

// A HashlimitEntty represents a single item in a bucket.
type HashlimitEntry struct {
	// The first field is the time to expiry in seconds.
	Expires uint64

	// The second field shows what packets are matched, as defined by the relevant `--hashlimit-mode` option.
	Match string

	// Current (renaming) credits for this entry.
	Current uint64

	// Maximum credits this entry can accrue.
	Capacity uint64

	// How many credits are deducted for every successful match.
	Cost uint64
}

// List the base path for buckets
func HashlimitBuckets() ([]string, error) {
	dir, err := os.ReadDir(HashlimitBasePath)
	if err != nil {
		return nil, err
	}
	buckets := make([]string, 0, len(dir))
	for _, d := range dir {
		buckets = append(buckets, d.Name())
	}
	return buckets, nil
}

func HashlimitListBucket(bucket string) ([]HashlimitEntry, error) {
	f, err := os.Open(filepath.Join(HashlimitBasePath, bucket))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	items := make([]HashlimitEntry, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) != 5 {
			continue
		}
		item := HashlimitEntry{Match: fields[1]}
		item.Expires, err = strconv.ParseUint(fields[0], 10, 64)
		if err != nil {
			return nil, err
		}
		item.Current, err = strconv.ParseUint(fields[2], 10, 64)
		if err != nil {
			return nil, err
		}
		item.Capacity, err = strconv.ParseUint(fields[3], 10, 64)
		if err != nil {
			return nil, err
		}
		item.Cost, err = strconv.ParseUint(fields[4], 10, 64)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}
