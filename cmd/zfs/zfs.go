package zfs

import (
	"github.com/iBug/ibugtool/pkg/util"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "zfs",
	Short: "ZFS commands",
	Args:  cobra.NoArgs,
	RunE:  util.ShowHelp,
}
