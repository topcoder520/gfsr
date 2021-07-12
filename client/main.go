package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">>> ")
		cmd, err := reader.ReadString('\n')
		fmt.Println(os.Args)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		cmd = strings.TrimSuffix(cmd, "\n")
		if strings.HasSuffix(cmd, "\r") {
			cmd = strings.TrimSuffix(cmd, "\r")
		}
		cmdArgs := strings.Split(cmd, " ")
		var cmdExe *exec.Cmd
		if len(cmdArgs) > 1 {
			cmdExe = exec.Command(cmdArgs[0], cmdArgs[1:]...)
		} else {
			cmdExe = exec.Command(cmdArgs[0])
		}

		cmdExe.Stderr = os.Stderr
		cmdExe.Stdout = os.Stdout
		err = cmdExe.Run()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
