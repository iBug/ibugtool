package cmd

import (
	"fmt"
	_ "unsafe"

	"github.com/spf13/cobra"
)

//go:linkname version main.version
var version string = "<unknown>"

var versionCmd = &cobra.Command{
	Use:    "version",
	Short:  "Print version and exit",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(cmd.OutOrStdout(), cmd.Root().Name(), "version", cmd.Root().Version)
	},
}
