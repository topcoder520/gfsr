package logs

import (
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func InitSysLog(dir string) {
	path, err := filepath.Abs(dir)
	if err != nil {
		ErrorF("Abs file path, sys.log err: ", err)
		return
	}
	os.MkdirAll(path, 0600)
	logFile, err := os.OpenFile(filepath.Clean(path+"/sys.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		ErrorF("open sys.log err: ", err)
		return
	}
	log.SetOutput(logFile)
}

func Println(msg interface{}) {
	log.Println(msg)
}

func Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func Error(err error) {
	log.Println(errors.WithStack(err))
}

func ErrorF(format string, err error) {
	log.Printf(format, errors.WithStack(err))
}

func ExitWithError(err error) {
	Error(err)
	os.Exit(1)
}
