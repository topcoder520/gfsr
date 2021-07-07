package middleware

import "net/http"

func MiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		next.ServeHTTP(rw, r)
	})
}
