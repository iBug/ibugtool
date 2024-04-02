package pve

import (
	"fmt"
	"strconv"

	"github.com/iBug/ibugtool/pkg/pve"
	"github.com/spf13/cobra"
)

var findCTCmd = &cobra.Command{
	Use:  "findct PID",
	RunE: fincCTRunE,
	Args: cobra.ExactArgs(1),
}

func init() {
	Cmd.AddCommand(findCTCmd)
}

func fincCTRunE(cmd *cobra.Command, args []string) error {
	pid, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	ctID, err := pve.FindCTByPid(pid)
	if err == pve.ErrProcessNotLXC {
		cmd.SilenceUsage = true
	}
	if err != nil {
		return err
	}
	fmt.Fprintln(cmd.OutOrStdout(), ctID)
	return nil
}
