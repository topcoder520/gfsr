package handler

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/topcoder520/gfsr/fileserver/config"
	"github.com/topcoder520/gfsr/fileserver/cryutil"
	"github.com/topcoder520/gfsr/fileserver/dao"
	"github.com/topcoder520/gfsr/fileserver/logs"
	"github.com/topcoder520/gfsr/fileserver/model"
	"github.com/topcoder520/gfsr/fileserver/session"
	"github.com/topcoder520/goutil"
)

func init() {
}

func GetServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	fileServerHandler := http.StripPrefix("/files/", http.FileServer(http.Dir(config.AbsDir)))

	mux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(rw, "ok")
	})

	mux.HandleFunc("/static/", func(rw http.ResponseWriter, r *http.Request) {
		staticResource, err := filepath.Abs(path.Join("./webapp", r.RequestURI))
		if err != nil {
			rw.WriteHeader(http.StatusNotFound)
			logs.Error(err)
			fmt.Fprintln(rw, "404 page not found")
		}
		http.ServeFile(rw, r, staticResource)
	})

	mux.HandleFunc("/login", func(rw http.ResponseWriter, r *http.Request) {
		loginPage, err := filepath.Abs("./webapp/view/login.html")
		if err != nil {
			rw.WriteHeader(http.StatusNotFound)
			logs.Error(err)
			fmt.Fprintln(rw, "404 page not found")
		}
		http.ServeFile(rw, r, loginPage)
	})

	mux.HandleFunc("/api/login", func(rw http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logs.Println(err)
			}
		}()
		if strings.ToLower(r.Method) == "post" {
			err := r.ParseForm()
			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(rw, "bad request")
				return
			}
			username := r.PostForm.Get("username")
			password := r.PostForm.Get("password")
			if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
				message := &model.Message{Code: 300, Msg: "please input username or password"}
				j, _ := json.Marshal(message)
				fmt.Fprintln(rw, string(j))
				return
			}
			pwd, _ := cryutil.EncryptDesCBC([]byte(username+password), []byte(config.PwdCryKey))
			password = fmt.Sprintf("%x", md5.Sum(pwd))
			if !dao.ExistUser(username, password) {
				message := &model.Message{Code: 300, Msg: "please input correct username or correct password"}
				j, _ := json.Marshal(message)
				fmt.Fprintln(rw, string(j))
				return
			}

			session := session.GetSession(rw, r)
			plainText := r.RemoteAddr + "@" + goutil.Parse_datetime_to_timestr(time.Now(), 1)
			token, err := cryutil.EncryptDesCBC([]byte(plainText), []byte(config.TokenKey))
			if err != nil {
				fmt.Fprintln(rw, http.StatusInternalServerError)
				return
			}
			session.Set(config.TokenName, token)
			session.Set("username", username)
			cookie := http.Cookie{
				Name:     config.TokenName,
				Value:    string(token),
				Path:     "/",
				HttpOnly: true,
				MaxAge:   int(config.TokenTimeOut),
			}
			http.SetCookie(rw, &cookie)

			message := &model.Message{Code: 200}
			j, err := json.Marshal(message)
			if err != nil {
				fmt.Fprintln(rw, http.StatusInternalServerError)
				return
			}
			fmt.Fprintln(rw, string(j))
		} else {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintln(rw, "request method do not allow 'GET'")
		}
	})

	mux.HandleFunc("/files/", func(rw http.ResponseWriter, r *http.Request) {
		fileServerHandler.ServeHTTP(rw, r)
	})
	return mux
}
