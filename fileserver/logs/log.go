package logs

import (
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func init() {
	path, err := filepath.Abs("./log/sys.log")
	if err != nil {
		ErrorF("Abs file path, sys.log err: ", err)
		return
	}
	logFile, err := os.Create(path)
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
