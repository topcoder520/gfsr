package handler

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/topcoder520/gfsr/fileserver/config"
	"github.com/topcoder520/gfsr/fileserver/logs"
)

func init() {
}

func GetServeMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(rw, "uklj")
	})

	mux.HandleFunc("/login", func(rw http.ResponseWriter, r *http.Request) {
		loginPage, err := filepath.Abs("./view/login.html")
		if err != nil {
			rw.WriteHeader(http.StatusNotFound)
			logs.Error(err)
			fmt.Fprintln(rw, "404 page not found")
		}
		http.ServeFile(rw, r, loginPage)
	})

	mux.HandleFunc("/api/login", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(rw, "uklj")
	})

	fileServerHandler := http.StripPrefix("/files/", http.FileServer(http.Dir(config.AbsDir)))
	mux.HandleFunc("/files/", func(rw http.ResponseWriter, r *http.Request) {
		fileServerHandler.ServeHTTP(rw, r)
	})
	return mux
}
