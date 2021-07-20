package handler

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
		//http.RedirectHandler("/files/", http.StatusOK)
	})

	mux.HandleFunc("/files/", func(rw http.ResponseWriter, r *http.Request) {
		fileServerHandler.ServeHTTP(rw, r)
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
			strToken := fmt.Sprintf("%x", md5.Sum(token))
			session.Set(config.TokenName, strToken)
			session.Set(strToken, plainText)
			session.Set("username", username)
			rw.Header().Set(config.TokenName, strToken)
			Success(nil, rw)
		} else {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintln(rw, "request method do not allow 'GET'")
		}
	})

	mux.HandleFunc("/api/files/", func(rw http.ResponseWriter, r *http.Request) {
		handler := http.StripPrefix("/api/files/", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			p := path.Join("/", path.Clean(r.URL.Path))
			f, err := os.Open(path.Join(config.AbsDir, p))
			if err != nil {
				Fail(500, err.Error(), rw)
				return
			}
			defer f.Close()
			fileInfos, err := f.Readdir(0)
			if err != nil {
				Fail(500, err.Error(), rw)
				return
			}
			rs := make([]model.FileInfo, len(fileInfos))
			for _, fileInfo := range fileInfos {
				file := model.FileInfo{
					Name:    fileInfo.Name(),
					Path:    path.Join(r.URL.Path, fileInfo.Name()),
					Size:    int(fileInfo.Size()),
					IsDir:   fileInfo.IsDir(),
					ModTime: fileInfo.ModTime().Format("2006-01-02 15:04:05"),
					Mode:    fileInfo.Mode().String(),
				}
				rs = append(rs, file)
			}
			Success(rs, rw)
		}))
		handler.ServeHTTP(rw, r)
	})

	mux.HandleFunc("/api/handlecmd/", func(rw http.ResponseWriter, r *http.Request) {
		handler := http.StripPrefix("/api/handlecmd/", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			cmd := r.URL.Path
			switch cmd {
			case "cd":
				cdpath := r.URL.Query().Get("cdpath")
				if len(strings.Trim(cdpath, " ")) == 0 {
					Fail(500, "cdpath is empty string", rw)
					return
				}
				cdpath = path.Join("/", path.Clean(cdpath))
				cdpath = path.Join(config.AbsDir, cdpath)
				f, err := os.Open(cdpath)
				if err != nil {
					Fail(int(ErrorFileNotFoundStatus), err.Error(), rw)
					return
				}
				fi, err := f.Stat()
				if err != nil {
					Fail(int(ErrorFileNotFoundStatus), err.Error(), rw)
					return
				}
				if !fi.IsDir() {
					Fail(int(ErrorFileNotDirStatus), "path is not direct", rw)
					return
				}
				Success("", rw)
				return
			}
			Success("", rw)
		}))
		handler.ServeHTTP(rw, r)
	})

	return mux
}
