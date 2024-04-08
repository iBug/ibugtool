package cmd

import (
	"fmt"
	_ "unsafe"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:    "version",
	Short:  "Print version and exit",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		root := cmd.Root()
		fmt.Fprintln(cmd.OutOrStdout(), root.Name(), "version", root.Version)
	},
}
