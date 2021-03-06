package main

import (
	"github.com/topcoder520/gfsr/fileserver/cmd"
	"github.com/topcoder520/gfsr/fileserver/logs"
)

func init() {
	logDir := "./log"
	logs.InitSysLog(logDir)
	logs.InitAccessLog(logDir)
}

func main() {
	if err := cmd.RootCommand.Execute(); err != nil {
		logs.ExitWithError(err)
	}
}
