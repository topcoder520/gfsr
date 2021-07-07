package middleware

import "net/http"

var limiter = NewIPRateLimiter(1, 10)

func FileServerMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		//访问日志
		//访问限速
		limiter := limiter.GetLimiter(r.RemoteAddr)
		if !limiter.Allow() {
			http.Error(rw, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}
		//登录验证

		next.ServeHTTP(rw, r)
	})
}
