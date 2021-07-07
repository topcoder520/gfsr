package handler

import (
	"fmt"
	"net/http"
)

var ServeMux *http.ServeMux

func init() {
	ServeMux = http.NewServeMux()
	routeHandler()
}

func routeHandler() {
	ServeMux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(rw, "uklj")
	})
}
