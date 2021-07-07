package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var ErrorExit = errors.New("exit os")

var RootCommand = &cobra.Command{
	Use:   "fileserver",
	Short: "file server",
}
