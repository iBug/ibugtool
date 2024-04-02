package cmd

import (
	"github.com/iBug/ibugtool/cmd/pve"
	"github.com/iBug/ibugtool/pkg/util"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:  "ibugtool",
	RunE: util.ShowHelp,
	CompletionOptions: cobra.CompletionOptions{
		HiddenDefaultCmd: true,
	},
	Version: version,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(pve.PveCmd)
	rootCmd.AddCommand(versionCmd)
}
