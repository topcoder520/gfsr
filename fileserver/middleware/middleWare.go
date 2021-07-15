package middleware

import (
	"net/http"
	"strings"

	"github.com/topcoder520/gfsr/fileserver/config"
	"github.com/topcoder520/gfsr/fileserver/logs"
	"github.com/topcoder520/gfsr/fileserver/session"
)

var limiter = NewIPRateLimiter(1, 10)

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
		logs.AccessLog.Printf("method:%s | ip:%s | referer:%s | userAgent:%s\n[END]\n",
			r.Method, r.RemoteAddr, r.Referer(), r.Header["User-Agent"])
		//登录验证
		uri := r.RequestURI
		if !strings.HasPrefix(uri, "/login") && !strings.HasPrefix(uri, "/api/login") && !strings.HasPrefix(uri, "/static/") {
			token := session.Get(config.TokenName)
			if token == nil {
				//跳转登录页面
				http.Redirect(rw, r, "/login", http.StatusFound)
				return
			}
			//验证token
			//获取请求的token
			reToken := r.Header.Get(config.TokenName)
			if reToken != token {
				//跳转登录页面
				http.Redirect(rw, r, "/login", http.StatusFound)
				return
			}
			strtoken := session.Get(token)
			fiels := strings.Split(strtoken.(string), "@")
			if len(fiels) < 2 {
				logs.Println("fiels len less 2")
				http.Redirect(rw, r, "/login", http.StatusFound)
				return
			}
			if r.RemoteAddr != fiels[0] {
				http.Redirect(rw, r, "/login", http.StatusFound)
				return
			}
		}
		next.ServeHTTP(rw, r)
	})
}
