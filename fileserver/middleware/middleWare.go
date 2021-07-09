package middleware

import (
	"net/http"

	"github.com/topcoder520/gfsr/fileserver/logs"
	"github.com/topcoder520/gfsr/fileserver/session"
)

var limiter = NewIPRateLimiter(1, 10)
var tokenName = "token_/1w#-=9*x8m,89f%^7&"

func FileServerMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		//访问限速
		limiter := limiter.GetLimiter(r.RemoteAddr)
		if !limiter.Allow() {
			http.Error(rw, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}
		//为每个请求生成session
		session := session.GetSession(rw, r)
		//访问日志
		r.ParseForm()
		logs.AccessLog.Printf("method:%s | ip:%s | referer:%s | userAgent:%s\n[END]\n",
			r.Method, r.RemoteAddr, r.Referer(), r.Header["User-Agent"])
		//登录验证
		uri := r.RequestURI
		if uri != "/login" {
			token := session.Get(tokenName)
			if token == nil {
				//跳转登录页面
				return
			}
			//获取请求的token
			if r.Header.Get("LOGIN-ACESS-TOKEN") != token {
				//跳转登录页面
				return
			}
		}
		next.ServeHTTP(rw, r)
	})
}
