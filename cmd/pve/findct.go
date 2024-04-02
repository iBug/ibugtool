package pve

import (
	"github.com/iBug/ibugtool/pkg/util"
	"github.com/spf13/cobra"
)

var findCTCmd = &cobra.Command{
	Use:  "findct",
	RunE: util.ShowHelp,
	Args: cobra.ExactArgs(1),
}

func init() {
	Cmd.AddCommand(findCTCmd)
}
