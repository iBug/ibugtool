package pve

import (
	"github.com/iBug/ibugtool/pkg/util"
	"github.com/spf13/cobra"
)

var PveCmd = &cobra.Command{
	Use:   "pve",
	Short: "Proxmox VE management commands",
	RunE:  util.ShowHelp,
}
