package cmd

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:  "ibugtool",
	RunE: RunE,
}

func RunE(cmd *cobra.Command, args []string) error {
	return nil
}

func Execute() error {
	return RootCmd.Execute()
}
