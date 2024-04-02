//go:build linux

package zfs

import (
	"fmt"
	"time"

	"github.com/iBug/ibugtool/pkg/util"
	"github.com/iBug/ibugtool/pkg/zfs"
	"github.com/spf13/cobra"
)

var arcstatsCmd = &cobra.Command{
	Use:     "arcstats",
	Aliases: []string{"arc"},
	Short:   "Monitor ZFS ARC statistics",
	Args:    cobra.NoArgs,
	RunE:    arcstatsRunE,
}

var arcstatsInterval time.Duration

func init() {
	flags := arcstatsCmd.Flags()
	flags.DurationVarP(&arcstatsInterval, "interval", "i", 1*time.Second, "Interval between updates")

	Cmd.AddCommand(arcstatsCmd)
}

func arcstatsRunE(cmd *cobra.Command, args []string) error {
	last, err := zfs.GetARCStats()
	if err != nil {
		return err
	}

	fmt.Printf("%8s  %8s  %5s  %5s%%  %7s  %6s  %5s%%\n",
		"req/s", "ARC Hit", "Miss", "Rate", "L2 Hit", "Miss", "Rate")
	defer fmt.Println()

	interval := arcstatsInterval
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		s, err := zfs.GetARCStats()
		if err != nil {
			return err
		}

		hits := s.Hits - last.Hits
		misses := s.Misses - last.Misses
		l2hits := s.L2Hits - last.L2Hits
		l2misses := s.L2Misses - last.L2Misses

		hitrate := float64(hits) / interval.Seconds()
		missrate := float64(misses) / interval.Seconds()
		l2hitrate := float64(l2hits) / interval.Seconds()
		l2missrate := float64(l2misses) / interval.Seconds()
		reqrate := float64(hits+misses) / interval.Seconds()
		hitratio, l2hitratio := 0.0, 0.0
		if hits+misses > 0 {
			hitratio = float64(hits) / float64(hits+misses) * 100.0
		}
		if l2hits+l2misses > 0 {
			l2hitratio = float64(l2hits) / float64(l2hits+l2misses) * 100.0
		}
		fmt.Printf("%s%8.1f  %8.1f  %5.1f  %5.1f%%  %7.1f  %6.1f  %5.1f%%", util.ResetLine,
			reqrate, hitrate, missrate, hitratio, l2hitrate, l2missrate, l2hitratio)
		last = s
	}
	return nil
}
