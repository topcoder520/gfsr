package config

var (
	Port           int
	Dir            string //file server root dir
	AbsDir         string
	TokenName      = "LOGIN-ACESS-TOKEN"
	TokenKey       = "yuioplkj"
	TokenTimeOut   = 3600
	PwdCryKey      = "hjknmcop"
	CACertPath     string //ca 证书
	ServerCertPath string //服务器证书
	ServerKeyPath  string //服务器密钥
)
