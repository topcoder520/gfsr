package request

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	netUrl "net/url"
	"strings"

	"github.com/topcoder520/gfsr/client/config"
	"github.com/topcoder520/gfsr/client/httpclient"
	"github.com/topcoder520/gfsr/client/model"
)

var IsLogout bool = true //是否退出

type Status int

const (
	OKStatus                Status = 200 //请求成功
	FailLoginStatus         Status = 300 //登录失败
	ErrorFileNotFoundStatus Status = 400 //文件错误
	ErrorRequestStatus      Status = 500 //请求错误
)

var TokenValue string
var SessionID string
var MaxAge int

const (
	TokenName     string = "Login-Acess-Token"
	SessionIDName string = "gosessionid"
)

func GetSessionStat(resp *http.Response) {
	for k, v := range resp.Header {
		if k == TokenName {
			TokenValue = v[0]
		}
	}
	cookies := resp.Cookies()
	for _, c := range cookies {
		if c.Name == SessionIDName {
			SessionID = c.Value
			MaxAge = c.MaxAge
		}
	}
}

func SetSessionStat(req *http.Request) {
	req.Header.Set(TokenName, TokenValue)
	cookie := &http.Cookie{
		Name:     SessionIDName,
		Value:    SessionID,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int(MaxAge),
	}
	req.AddCookie(cookie)
}

func Get(uri string) ([]byte, error) {
	client := httpclient.OnceTLS()
	url := config.GetFileServerAddress() + uri
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	SetSessionStat(req)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	GetSessionStat(resp)
	rs, err := ioutil.ReadAll(resp.Body)
	return rs, err
}

func Post(uri string, formParam map[string]string) ([]byte, error) {
	client := httpclient.OnceTLS()
	url := config.GetFileServerAddress() + uri
	value := make(netUrl.Values)
	for k, v := range formParam {
		value.Add(k, v)
	}
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(value.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	SetSessionStat(req)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	GetSessionStat(resp)
	rs, err := ioutil.ReadAll(resp.Body)
	return rs, err
}

func ParseResult(rs []byte) (*model.Message, error) {
	message := &model.Message{}
	err := json.Unmarshal(rs, message)
	if err != nil {
		return nil, err
	}
	CheckSessionTimeOut(message)
	return message, nil
}

type FormData map[string]string

func PostData() FormData {
	return FormData{}
}

func (data FormData) Add(key, value string) {
	data[key] = value
}

func CheckSessionTimeOut(message *model.Message) {
	if message.Code == int(FailLoginStatus) {
		IsLogout = true
	} else if message.Code == int(OKStatus) {
		IsLogout = false
	}
}
