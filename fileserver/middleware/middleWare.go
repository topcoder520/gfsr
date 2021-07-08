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
		logs.AccessLog.Printf("method:%s\nip:%s\nreferer:%s\nurl:%s\nuserAgent:%s\n",
			r.Method, r.RemoteAddr, r.Referer(), r.URL, r.Header["User-Agent"])
		logs.AccessLog.Println("header:", r.Header)
		logs.AccessLog.Println("formInfo:", r.PostForm)
		logs.AccessLog.Println("[END]")
		//登录验证
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
		next.ServeHTTP(rw, r)
	})
}
