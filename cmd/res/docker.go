package res

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/iBug/ibugtool/pkg/docker"
	"github.com/iBug/ibugtool/pkg/proc"
	"github.com/iBug/ibugtool/pkg/util"
	"github.com/olekukonko/tablewriter"
	"github.com/prometheus/procfs"
	"github.com/spf13/cobra"
)

var dockerCmd = &cobra.Command{
	Use:   "docker <container...>",
	Short: "Display Docker container resource usage",
	Args:  cobra.MinimumNArgs(1),
	RunE:  dockerRunE,
}

var (
	dockerSortBy string
	dockerTopK   int
)

func dockerRunE(cmd *cobra.Command, args []string) error {
	now := time.Now()
	now_ts := float64(now.Unix()) + float64(now.UnixNano())/1e9
	cpuRatio := func(stat procfs.ProcStat) float64 {
		startTime, err := stat.StartTime()
		if err != nil {
			return 0
		}
		return stat.CPUTime() / (now_ts - startTime)
	}

	var sortFunc func(a, b procfs.Proc) int
	switch dockerSortBy {
	case "memory":
		sortFunc = func(a, b procfs.Proc) int {
			as, err := a.Stat()
			if err != nil {
				return 0
			}
			bs, err := b.Stat()
			if err != nil {
				return 0
			}
			return int(bs.ResidentMemory() - as.ResidentMemory())
		}
	case "cpu":
		sortFunc = func(a, b procfs.Proc) int {
			as, err := a.Stat()
			if err != nil {
				return 0
			}
			bs, err := b.Stat()
			if err != nil {
				return 0
			}
			diff := cpuRatio(bs) - cpuRatio(as)
			switch {
			case diff < 0:
				return -1
			case diff > 0:
				return 1
			default:
				return 0
			}
		}
	default:
		return fmt.Errorf("invalid sort option %q", dockerSortBy)
	}

	ctx := cmd.Context()
	for _, container := range args {
		pids, err := docker.ContainerProcesses(ctx, container)
		if err != nil {
			return fmt.Errorf("could not get processes info for container %s: %v", container, err)
		}

		procs := make([]procfs.Proc, 0, len(pids))
		for _, pid := range pids {
			p, err := proc.FS.Proc(int(pid))
			if err != nil {
				return fmt.Errorf("could not open process %d: %v", pid, err)
			}
			procs = append(procs, p)
		}
		slices.SortFunc(procs, sortFunc)

		// Output
		table := util.DefaultTable(cmd.OutOrStdout())
		tAlignment := []int{
			tablewriter.ALIGN_RIGHT,
			tablewriter.ALIGN_RIGHT,
			tablewriter.ALIGN_RIGHT,
			tablewriter.ALIGN_DEFAULT,
		}
		tHeaders := []string{"PID", "CPU%", "Memory", "Command"}
		table.SetAutoWrapText(false)
		table.SetColumnAlignment(tAlignment)
		table.Append(tHeaders)

		for _, p := range procs {
			stat, err := p.Stat()
			if err != nil {
				return fmt.Errorf("could not stat process %d: %v", p.PID, err)
			}
			cpu := fmt.Sprintf("%.1f%%", cpuRatio(stat)*100)
			mem := util.FormatSize(uint64(stat.ResidentMemory()))
			cmdline, err := p.CmdLine()
			if err != nil {
				return fmt.Errorf("could not get command line for process %d: %v", p.PID, err)
			}
			table.Append([]string{
				fmt.Sprintf("%d", p.PID),
				cpu,
				mem,
				strings.Join(cmdline, " "),
			})
		}
		table.Render()
	}

	// Implement Docker container resource usage display
	return nil
}

func init() {
	flags := dockerCmd.Flags()
	flags.StringVarP(&dockerSortBy, "sort", "s", "memory", "Sort by (memory or cpu)")
	flags.IntVarP(&dockerTopK, "topk", "k", 0, "Show only the top K processes (0 means all)")
	Cmd.AddCommand(dockerCmd)
}
