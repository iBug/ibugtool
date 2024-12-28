package res

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/iBug/ibugtool/pkg/docker"
	"github.com/iBug/ibugtool/pkg/proc"
	"github.com/iBug/ibugtool/pkg/util"
	"github.com/olekukonko/tablewriter"
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
	var sortFunc func(a, b proc.ProcInfo) int
	switch dockerSortBy {
	case "mem", "memory":
		sortFunc = func(a, b proc.ProcInfo) int {
			return int(b.ResidentMemory - a.ResidentMemory)
		}
	case "cpu":
		sortFunc = func(a, b proc.ProcInfo) int {
			diff := b.CPURatio() - a.CPURatio()
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
	w := cmd.OutOrStdout()
	for _, containerID := range args {
		container, err := docker.Client().ContainerInspect(ctx, containerID)
		if err != nil {
			return fmt.Errorf("could not inspect container %s: %v", containerID, err)
		}
		fmt.Fprintf(w, "Container ID: %s\n", container.ID[:12])
		fmt.Fprintf(w, "Container Name: %s\n", strings.TrimPrefix(container.Name, "/"))

		pids, err := docker.ContainerProcesses(ctx, containerID)
		if err != nil {
			return fmt.Errorf("could not get processes info for container %s: %v", containerID, err)
		}

		infos := make([]proc.ProcInfo, 0, len(pids))
		for _, pid := range pids {
			info, err := proc.GetProcInfo(pid)
			if err != nil {
				return fmt.Errorf("could not open process %d: %v", pid, err)
			}
			infos = append(infos, info)
		}
		slices.SortFunc(infos, sortFunc)

		// Output
		table := util.DefaultTable(w)
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

		for k, info := range infos {
			if dockerTopK > 0 && k >= dockerTopK {
				break
			}
			cpu := fmt.Sprintf("%.1f%%", info.CPURatio()*100)
			mem := util.FormatSize(info.ResidentMemory)
			table.Append([]string{
				strconv.Itoa(info.Pid),
				cpu,
				mem,
				strings.Join(info.Cmdline, " "),
			})
		}
		table.Render()
		w.Write([]byte("\n"))
	}
	return nil
}

func init() {
	flags := dockerCmd.Flags()
	flags.StringVarP(&dockerSortBy, "sort", "s", "memory", "Sort by (memory or cpu)")
	flags.IntVarP(&dockerTopK, "topk", "k", 0, "Show only the top K processes (0 means all)")
	Cmd.AddCommand(dockerCmd)
}
