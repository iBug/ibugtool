package cmd

import (
	"github.com/iBug/ibugtool/cmd/net"
	"github.com/iBug/ibugtool/cmd/pve"
	"github.com/iBug/ibugtool/cmd/res"
	"github.com/iBug/ibugtool/cmd/zfs"
	"github.com/iBug/ibugtool/internal/version"
	"github.com/iBug/ibugtool/pkg/util"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:  "ibugtool",
	Args: cobra.NoArgs,
	RunE: util.ShowHelp,
	CompletionOptions: cobra.CompletionOptions{
		HiddenDefaultCmd: true,
	},
	Version: version.Version,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.AddCommand(
		updateCmd,
		versionCmd,
		// Entrypoint command for sub-packages are always "Cmd" or "MakeCmd()"
		res.Cmd,
		net.HashlimitCmd,
		pve.Cmd,
		zfs.Cmd,
	)
}
