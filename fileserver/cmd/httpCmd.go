package cmd

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/topcoder520/gfsr/fileserver/config"
	"github.com/topcoder520/gfsr/fileserver/handler"
	"github.com/topcoder520/gfsr/fileserver/logs"
	"github.com/topcoder520/gfsr/fileserver/middleware"
)

var httpCommand = &cobra.Command{
	Use:     "http",
	Short:   "HTTP protocol file server",
	Example: "gofileserver http [-p 80] [-d path] ",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := startServer(args); err != nil {
			logs.Error(err)
			return err
		}
		return nil
	},
}

func init() {
	RootCommand.AddCommand(httpCommand)
	httpCommand.Flags().IntVarP(&config.Port, "port", "p", 80, "The port of server listening on")
	httpCommand.Flags().StringVarP(&config.Dir, "dir", "d", ".", "server directory")
}

func startServer(args []string) error {
	absPath, err := filepath.Abs(filepath.Clean(config.Dir))
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
	//http.FileServer(http.Dir(absPath))
	config.AbsDir = absPath
	logs.Printf("listening port: %d \n", config.Port)
	logs.Printf("server directory: %s\n", config.AbsDir)
	logs.Println("http fileserver started successfully.")
	return http.ListenAndServe(fmt.Sprintf(":%d", config.Port),
		middleware.FileServerMiddleWare(handler.GetServeMux()))
}
