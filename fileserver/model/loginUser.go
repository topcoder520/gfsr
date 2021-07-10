package model

type LoginUser struct {
	Id         int    `json:"id" col:"Id"`
	Username   string `json:"username" col:"Username"`
	Password   string `json:"-" col:"Password"`
	CreateTime string `json:"createtime" col:"createtime"`
}

type LoginLog struct {
	Id          int    `json:"id" col:"Id"`
	LoginUserId int    `json:"LoginUser" col:"LoginUserId"`
	CreateTime  string `json:"createtime" col:"CreateTime"`
}

type Message struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
