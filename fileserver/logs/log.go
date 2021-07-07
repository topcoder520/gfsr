package logs

import (
	"log"
	"os"

	"github.com/pkg/errors"
)

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
