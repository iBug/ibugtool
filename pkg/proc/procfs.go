package proc

import (
	"time"

	"github.com/prometheus/procfs"
)

const Root = procfs.DefaultMountPoint

var FS procfs.FS

func init() {
	var err error
	FS, err = procfs.NewDefaultFS()
	if err != nil {
		panic(err)
	}
}

type ProcInfo struct {
	Pid     int
	Comm    string
	Cmdline []string

	CPUTime        float64
	StartTime      float64
	ResidentMemory uint64

	queryTime float64
}

func GetProcInfo(pid int) (ProcInfo, error) {
	info := ProcInfo{Pid: pid}
	now := time.Now()
	info.queryTime = float64(now.Unix()) + float64(now.UnixNano())/1e9

	proc, err := FS.Proc(pid)
	if err != nil {
		return info, err
	}
	info.Comm, err = proc.Comm()
	if err != nil {
		return info, err
	}
	info.Cmdline, err = proc.CmdLine()
	if err != nil {
		return info, err
	}
	stat, err := proc.NewStat()
	if err != nil {
		return info, err
	}
	info.CPUTime = stat.CPUTime()
	info.StartTime, err = stat.StartTime()
	info.ResidentMemory = uint64(stat.ResidentMemory())
	if err != nil {
		return info, err
	}
	return info, nil
}

func (info ProcInfo) CPURatio() float64 {
	return info.CPUTime / (info.queryTime - info.StartTime)
}
