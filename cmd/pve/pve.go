package pve

import (
	"github.com/iBug/ibugtool/pkg/util"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "pve",
	Short: "Proxmox VE commands",
	Args:  cobra.NoArgs,
	RunE:  util.ShowHelp,
}
