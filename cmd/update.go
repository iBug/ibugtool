package cmd

import (
	"github.com/iBug/ibugtool/pkg/updater"
	"github.com/spf13/cobra"
)

var forceUpdate bool

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update this tool from GitHub releases",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		return updater.UpdateBinary(cmd.OutOrStdout(), forceUpdate)
	},
}

func init() {
	flags := updateCmd.Flags()
	flags.BoolVarP(&forceUpdate, "force", "f", false, "Force update")
}
