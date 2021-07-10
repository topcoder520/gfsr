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
			/* //验证token
			//获取请求的token
			if r.Header.Get("LOGIN-ACESS-TOKEN") != token {
				//跳转登录页面
				http.Redirect(rw, r, "/login", http.StatusFound)
				return
			}
			reqToken := r.Header.Get("LOGIN-ACESS-TOKEN")
			deToken, err := cryutil.DecryptDesCBC([]byte(reqToken), []byte(config.TokenKey))
			if err != nil {
				logs.Error(err)
				http.Redirect(rw, r, "/login", http.StatusFound)
				return
			}
			fiels := strings.Split(string(deToken), "@")
			if len(fiels) < 2 {
				logs.Println("fiels len less 2")
				http.Redirect(rw, r, "/login", http.StatusFound)
				return
			}
			if r.RemoteAddr != fiels[0] {
				http.Redirect(rw, r, "/login", http.StatusFound)
				return
			} */
		}
		next.ServeHTTP(rw, r)
	})
}
