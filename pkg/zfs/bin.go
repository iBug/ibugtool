package zfs

import "os/exec"

var BinPath = "/sbin/zfs"

func ZfsCommand(args ...string) *exec.Cmd {
	return exec.Command(BinPath, args...)
}
