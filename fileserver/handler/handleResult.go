package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/topcoder520/gfsr/fileserver/model"
)

type Status int

const (
	OKStatus                Status = 200 //请求成功
	FailLoginStatus         Status = 300 //登录失败
	ErrorFileNotFoundStatus Status = 400 //文件错误
	ErrorFileNotDirStatus   Status = 410 //文件不是目录错误
	ErrorRequestStatus      Status = 500 //请求错误
)

func Success(data interface{}, w http.ResponseWriter) {
	message := &model.Message{}
	message.Code = 200
	message.Data = data
	j, err := json.Marshal(message)
	if err != nil {
		fmt.Fprintln(w, http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, string(j))
}

func Fail(Code int, Msg string, w http.ResponseWriter) {
	message := &model.Message{}
	message.Code = Code
	message.Msg = Msg
	j, err := json.Marshal(message)
	if err != nil {
		fmt.Fprintln(w, http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, string(j))
}
