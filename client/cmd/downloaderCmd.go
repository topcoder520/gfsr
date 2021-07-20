package cmd

import (
	"fmt"
	"net/http"
	"path"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/topcoder520/gfsr/client/colorp"
	"github.com/topcoder520/gfsr/client/config"
	"github.com/topcoder520/gfsr/client/httpclient"
)

func InitDownloaderCmd() *cobra.Command {
	fdlr := &cobra.Command{
		Use:     "fdlr",
		Short:   "fast downloader",
		Example: "fdlr [-c=3] [-o=newfilename] filename",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			defer func() {
				if err := recover(); err != nil {
					colorp.ErrorPrintln(err.(string))
				}
			}()
			download(args)
		},
	}
	fdlr.Flags().IntVarP(&c, "count", "c", runtime.NumCPU(), "goroutines count default runtime.NumCPU()")
	fdlr.Flags().StringVarP(&outputFilename, "output", "o", "", "outout filename")
	return fdlr
}

var c int //down goroutines
var outputFilename string

func download(args []string) {
	file := path.Clean(args[0])
	filePath := path.Join(config.CurrentPath, file)
	filePath = path.Join("/files/", filePath)
	client := httpclient.OnceTLS()
	resp, err := client.Head(config.GetFileServerAddress() + filePath)
	if err != nil {
		colorp.ErrorPrintln(err.Error())
		return
	}
	defer resp.Body.Close()
	isChunk := false
	if resp.StatusCode == http.StatusOK && resp.Header.Get("Accept-Ranges") == "bytes" {
		isChunk = true
	}
	fileSize := resp.Header.Get("Content-Length")
	fmt.Println(isChunk, fileSize)

}

func SingleDownloader() {

}

func MultDownloader() {}
