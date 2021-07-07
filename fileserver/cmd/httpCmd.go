package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/topcoder520/gfsr/fileserver/logs"
)

var httpCommand = &cobra.Command{
	Use:     "http",
	Short:   "HTTP protocol file server",
	Example: "gofileserver http [-port 80] [-dir path] ",
	RunE: func(cmd *cobra.Command, args []string) error {
		return startServer(args)
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

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(rw, "uklj")
	})
	http.ListenAndServe(fmt.Sprintf(":%d", port), middleWare(mux))

	return nil
}

func middleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		logs.Println(r.URL.Query())
		if len(r.URL.Query().Get("skip")) > 0 {
			fmt.Fprintln(rw, "ok")
			return
		}
		next.ServeHTTP(rw, r)
	})
}
