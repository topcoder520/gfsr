package cmd

import (
	"path"

	"github.com/spf13/cobra"
	"github.com/topcoder520/gfsr/client/colorp"
	"github.com/topcoder520/gfsr/client/config"
	"github.com/topcoder520/gfsr/client/request"
)

func InitCdCmd() *cobra.Command {
	cdCmd := &cobra.Command{
		Use:   "cd",
		Short: "change directory",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			defer func() {
				if err := recover(); err != nil {
					colorp.ErrorPrintln(err.(string))
				}
			}()
			cdRunFunc(cmd, args)
		},
	}
	return cdCmd
}

func cdRunFunc(cmd *cobra.Command, args []string) {
	if args[0] == "~" {
		config.CurrentPath = "/"
		return
	}
	p := path.Join(config.CurrentPath, path.Clean(args[0]))
	message, err := request.HandleCdCmd(cmd.Use, p)
	if err != nil {
		colorp.ErrorPrintln(err.Error())
		return
	}
	if message.Code != int(request.OKStatus) {
		colorp.ErrorPrintln(message.Msg)
		return
	}
	config.CurrentPath = p
}
