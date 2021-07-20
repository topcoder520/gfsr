package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var RootCmd *cobra.Command

func NewRootCmd() {
	RootCmd = &cobra.Command{
		Use:   os.Args[0],
		Short: "file client",
	}
	RootCmd.AddCommand(InitLsCmd())
	RootCmd.AddCommand(InitPwdCmd())
	RootCmd.AddCommand(InitCdCmd())
	RootCmd.AddCommand(InitDownloaderCmd())
}

func ResetCommand() {
	if RootCmd != nil {
		RootCmd.ResetCommands()
		RootCmd = nil
	}
}
