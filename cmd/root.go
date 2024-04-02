package cmd

import (
	"github.com/iBug/ibugtool/pkg/util"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:  "ibugtool",
	RunE: util.ShowHelp,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(PveCmd)
}
