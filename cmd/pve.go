package cmd

import "github.com/spf13/cobra"

var PveCmd = &cobra.Command{
	Use:  "pve",
	RunE: PveRunE,
}

func PveRunE(cmd *cobra.Command, args []string) error {
	return cmd.Help()
}

var PveFindCTCmd = &cobra.Command{
	Use:  "findct",
	RunE: PveRunE,
	Args: cobra.ExactArgs(1),
}
