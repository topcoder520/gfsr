package cmd

import (
	"github.com/spf13/cobra"
	"github.com/topcoder520/gfsr/client/colorp"
	"github.com/topcoder520/gfsr/client/config"
)

func InitPwdCmd() *cobra.Command {
	PwdCmd := &cobra.Command{
		Use:   "pwd",
		Short: "print current path",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			defer func() {
				if err := recover(); err != nil {
					colorp.ErrorPrintln(err.(string))
				}
			}()
			colorp.WhitePrintln(config.CurrentPath)
		},
	}
	return PwdCmd
}
