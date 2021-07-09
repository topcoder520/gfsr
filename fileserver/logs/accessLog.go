package logs

import (
	"log"
	"os"
	"path/filepath"
)

var AccessLog *log.Logger

func InitAccessLog(dir string) {
	path, err := filepath.Abs(dir)
	if err != nil {
		ErrorF("Abs file path, access.log err: ", err)
		return
	}
	os.MkdirAll(path, 0600)

	logFile, err := os.OpenFile(filepath.Clean(path+"/access.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		ErrorF("open access.log err: ", err)
		return
	}
	AccessLog = log.New(logFile, "[INFO]", log.Ldate|log.Lmicroseconds)
}
