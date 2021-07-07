package cmd

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/topcoder520/gfsr/fileserver/logs"
	"github.com/topcoder520/gfsr/fileserver/middleware"
)

var httpCommand = &cobra.Command{
	Use:     "http",
	Short:   "HTTP protocol file server",
	Example: "gofileserver http [-port 80] [-dir path] ",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := startServer(args); err != nil {
			logs.Error(err)
			return err
		}
		return nil
	},
}

var port int
var dir string //file server root dir

func init() {
	RootCommand.AddCommand(httpCommand)
	httpCommand.Flags().IntVarP(&port, "port", "p", 80, "The port of server listening on")
	httpCommand.Flags().StringVarP(&dir, "dir", "d", ".", "server directory")
}

func startServer(args []string) error {
	absPath, err := filepath.Abs(filepath.Clean(dir))
	if err != nil {
		return err
	}
	f, err := os.Open(absPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		err = os.Mkdir(absPath, 0600)
		if err != nil {
			return err
		}
	} else {
		fileInfo, err := f.Stat()
		f.Close()
		if err != nil {
			return err
		}
		if !fileInfo.IsDir() {
			return errors.New("dir not directory")
		}
	}
	logs.Printf("listening port: %d \n", port)
	logs.Printf("server directory: %s\n", absPath)
	logs.Println("http fileserver started successfully.")
	return http.ListenAndServe(fmt.Sprintf(":%d", port), middleware.FileServerMiddleWare(http.FileServer(http.Dir(absPath))))
}
