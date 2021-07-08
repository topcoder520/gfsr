package logs

import (
	"log"
	"os"
	"path/filepath"
)

var AccessLog *log.Logger

func init() {
	path, err := filepath.Abs("./log/access.log")
	if err != nil {
		ErrorF("Abs file path, access.log err: ", err)
		return
	}
	logFile, err := os.Create(path)
	if err != nil {
		ErrorF("open access.log err: ", err)
		return
	}
	AccessLog = log.New(logFile, "[request]", log.Ldate|log.Lmicroseconds)
}
