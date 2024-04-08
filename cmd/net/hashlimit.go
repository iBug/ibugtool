package net

import (
	"fmt"

	"github.com/iBug/ibugtool/pkg/proc"
	"github.com/spf13/cobra"
)

var hashlimitCmd = &cobra.Command{
	Use:   "hashlimit [BUCKET]",
	Short: "Show information about the iptables `hashlimit` module",
	Long: `Show information about the iptables hashlimit module.

If no argument is given, lists all buckets.
If BUCKET is given, list entries in that bucket.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return hashlimitList(cmd)
		}
		return nil
	},
}

func hashlimitList(cmd *cobra.Command) error {
	out := cmd.OutOrStdout()
	buckets, err := proc.HashlimitBuckets()
	if err != nil {
		return err
	}
	for _, bucket := range buckets {
		fmt.Fprintln(out, bucket)
	}
	return nil
}

// maybe group this under a "net" command group
var HashlimitCmd = hashlimitCmd
