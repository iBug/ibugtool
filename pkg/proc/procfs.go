package proc

import "github.com/prometheus/procfs"

const Root = procfs.DefaultMountPoint

var FS procfs.FS

func init() {
	var err error
	FS, err = procfs.NewDefaultFS()
	if err != nil {
		panic(err)
	}
}
