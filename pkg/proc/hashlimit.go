package proc

import "os"

const HashlimitBasePath = "/proc/net/ipt_hashlimit"

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
