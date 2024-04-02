package util

import "github.com/spf13/cobra"

func ShowHelp(cmd *cobra.Command, args []string) error {
	return cmd.Help()
}
