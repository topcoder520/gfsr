package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/topcoder520/gfsr/fileserver/model"
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
