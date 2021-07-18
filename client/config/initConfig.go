package config

import "strings"

var (
	Protocol       string
	FileServerUrl  string
	FileServerPort string // :port
)

func init() {
	Protocol = "https"
	//FileServerUrl = "222.84.253.64"
	FileServerUrl = "localhost"
	//FileServerPort = ":8088"
	FileServerPort = ":8080"
}

func GetFileServerAddress() string {
	return strings.Join([]string{Protocol, "://", FileServerUrl, FileServerPort}, "")
}
