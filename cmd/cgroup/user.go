package cgroup

import (
	"fmt"
	osuser "os/user"
	"path/filepath"
	"slices"
	"strconv"

	"github.com/iBug/ibugtool/pkg/cgroupfs"
	"github.com/iBug/ibugtool/pkg/proc"
	"github.com/iBug/ibugtool/pkg/util"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "cgroup",
	Short: "Commands for managing cgroup information",
	RunE:  util.ShowHelp,
}

type uidInfo struct {
	Uid  uint64
	Info cgroupfs.CgroupInfo
}

func userMemShow(cmd *cobra.Command, args []string) error {
	meminfo, err := proc.GetMeminfo()
	if err != nil {
		return err
	}

	userSlices, err := cgroupfs.UserSlices()
	if err != nil {
		return err
	}
	infos := make([]uidInfo, 0, len(userSlices))
	for _, slice := range userSlices {
		info := uidInfo{}
		_, err := fmt.Sscanf(slice, "user-%d.slice", &info.Uid)
		if err != nil {
			continue
		}

		info.Info, err = cgroupfs.GetCgroupInfo(filepath.Join(cgroupfs.CgroupUserRoot, slice))
		if err != nil {
			return err
		}
		infos = append(infos, info)
	}
	slices.SortFunc(infos, func(a, b uidInfo) int {
		diff := int64(b.Info.MemoryCurrent - a.Info.MemoryCurrent)
		switch {
		case diff < 0:
			return -1
		case diff > 0:
			return 1
		default:
			return 0
		}
	})

	table := util.DefaultTable(cmd.OutOrStdout())
	tAlignment := []int{
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_DEFAULT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_RIGHT,
	}
	tHeaders := []string{"UID", "Name", "Memory", "Memory%", "Swap", "PIDs"}
	table.SetColumnAlignment(tAlignment)
	table.SetHeader(tHeaders)
	for _, info := range infos {
		var username string
		user, err := osuser.LookupId(strconv.FormatUint(info.Uid, 10))
		if err != nil {
			username = fmt.Sprintf("(%d)", info.Uid)
		} else {
			username = user.Username
		}
		row := []string{
			strconv.FormatUint(info.Uid, 10),
			username,
			util.FormatSizeAligned(info.Info.MemoryCurrent),
			fmt.Sprintf("%.1f %%", float64(info.Info.MemoryCurrent)/float64(meminfo.MemTotal)*100),
			util.FormatSizeAligned(info.Info.MemorySwapCurrent),
			fmt.Sprintf("%d", info.Info.Pids),
		}
		table.Append(row)
	}
	table.Render()
	return nil
}

var userMemCmd = &cobra.Command{
	Use:   "usermem",
	Short: "Show user memory usage",
	Args:  cobra.NoArgs,
	RunE:  userMemShow,
}

func init() {
	Cmd.AddCommand(userMemCmd)
}
