package main

import (
	"bufio"
	"os"
	"strings"

	"github.com/topcoder520/gfsr/client/cmd"
	"github.com/topcoder520/gfsr/client/colorp"
	"github.com/topcoder520/gfsr/client/request"
	"github.com/topcoder520/gfsr/client/sha1"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		if request.IsLogout {
			err := Start()
			if err != nil {
				colorp.ErrorPrintln(err.Error())
				continue
			}
		}
		colorp.LineHeaderInfo()
		args, err := reader.ReadString('\n')
		if err != nil {
			colorp.ErrorPrintln(err.Error())
			continue
		}
		args = strings.TrimSuffix(args, "\n")
		if strings.HasSuffix(args, "\r") {
			args = strings.TrimSuffix(args, "\r")
		}
		cmdArgs := strings.Split(args, " ")
		os.Args = os.Args[0:1]
		os.Args = append(os.Args, cmdArgs...)
		cmd.NewRootCmd()
		if err = cmd.RootCmd.Execute(); err != nil {
			colorp.ErrorPrintln(err.Error())
			continue
		}
		cmd.ResetCommand()
	}
}

func Start() error {
	/* bufReader := bufio.NewReader(os.Stdin)
	fmt.Print("please input username: ")
	username, err := bufReader.ReadString('\n')
	if err != nil {
		return err
	}
	username = strings.Trim(username, "\r\n")
	fmt.Print("please input password: ")
	pwd, err := gopass.GetPasswd()
	if err != nil {
		return err
	}
	pwdStr := strings.Trim(string(pwd), " ") */
	//登录操作
	//err = request.Login(username, sha1.Hex_sha1(pwdStr))
	err := request.Login("huangjing", sha1.Hex_sha1("huangjing123511"))
	if err != nil {
		return err
	}
	return nil
}
