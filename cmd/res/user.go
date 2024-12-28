package res

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
	Use:   "res",
	Short: "Commands for displaying system resource information",
	Args:  cobra.NoArgs,
	RunE:  util.ShowHelp,
}

type uidInfo struct {
	Uid  uint64
	Info cgroupfs.CgroupInfo
}

func userRssSize(uid uint64) (uint64, error) {
	userSlice := filepath.Join(cgroupfs.CgroupUserRoot, fmt.Sprintf("user-%d.slice", uid))
	pids, err := cgroupfs.GetAllPids(userSlice)
	if err != nil {
		return 0, err
	}
	var rss uint64
	for _, pid := range pids {
		p, err := proc.FS.Proc(int(pid))
		if err != nil {
			continue
		}
		// status, err := p.NewStatus()
		stat, err := p.Stat()
		if err != nil {
			continue
		}
		// rss += status.RssAnon + status.RssFile + status.RssShmem
		rss += uint64(stat.ResidentMemory())
	}
	return rss, nil
}

func userMemShow(cmd *cobra.Command, args []string) error {
	meminfo, err := proc.FS.Meminfo()
	if err != nil {
		return err
	}
	memTotal := 1024 * *meminfo.MemTotal

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
	}
	tHeaders := []string{"UID", "Name", "Memory", "Memory%"}
	if userMemIncludeRss {
		// The two columns for RSS
		tAlignment = append(tAlignment, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_RIGHT)
		tHeaders = append(tHeaders, "RSS", "RSS%")
	}
	tAlignment = append(tAlignment, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_RIGHT)
	tHeaders = append(tHeaders, "Swap", "PIDs")
	table.SetColumnAlignment(tAlignment)
	table.Append(tHeaders)

	var memorySum, memorySwapSum, pidsSum, rssSum uint64
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
			fmt.Sprintf("%.1f %%", float64(info.Info.MemoryCurrent)/float64(memTotal)*100),
		}
		if userMemIncludeRss {
			rss, err := userRssSize(info.Uid)
			if err != nil {
				cmd.PrintErrf("Failed to get RSS for UID %d: %v\n", info.Uid, err)
			}
			row = append(row,
				util.FormatSizeAligned(rss),
				fmt.Sprintf("%.1f %%", float64(rss)/float64(memTotal)*100))
			rssSum += rss
		}
		row = append(row,
			util.FormatSizeAligned(info.Info.MemorySwapCurrent),
			fmt.Sprintf("%d", info.Info.Pids))
		table.Append(row)

		memorySum += info.Info.MemoryCurrent
		memorySwapSum += info.Info.MemorySwapCurrent
		pidsSum += info.Info.Pids
	}
	footerLine := []string{
		"", "Total",
		util.FormatSizeAligned(memorySum),
		fmt.Sprintf("%.1f %%", float64(memorySum)/float64(memTotal)*100),
	}
	if userMemIncludeRss {
		footerLine = append(footerLine,
			util.FormatSizeAligned(rssSum),
			fmt.Sprintf("%.1f %%", float64(rssSum)/float64(memTotal)*100))
	}
	footerLine = append(footerLine,
		util.FormatSizeAligned(memorySwapSum),
		fmt.Sprintf("%d", pidsSum))
	table.Append(footerLine)
	table.Render()
	return nil
}

var userMemCmd = &cobra.Command{
	Use:   "usermem",
	Short: "Show user memory usage",
	Args:  cobra.NoArgs,
	RunE:  userMemShow,
}

var userMemIncludeRss bool

func init() {
	userMemCmd.Flags().BoolVarP(&userMemIncludeRss, "rss", "r", false, "Include resident set size")
	Cmd.AddCommand(userMemCmd)
}
