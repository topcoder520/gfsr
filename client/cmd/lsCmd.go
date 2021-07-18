package cmd

import (
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/topcoder520/gfsr/client/config"
	"github.com/topcoder520/gfsr/client/model"
	"github.com/topcoder520/gfsr/client/request"
)

var l bool

func InitLsCmd() *cobra.Command {
	LsCmd := &cobra.Command{
		Use:   "ls",
		Short: "list files",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(err)
				}
			}()
			RunFunc(args)
		},
	}
	LsCmd.Flags().BoolVarP(&l, "l", "l", false, "list file info ")
	return LsCmd
}

func RunFunc(args []string) {
	lsPath := ""
	if len(args) > 0 {
		lsPath = path.Join(config.CurrentPath, path.Clean(args[0]))
	} else {
		lsPath = path.Join(config.CurrentPath)
	}
	message, err := request.GetDirFiles(lsPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	if message.Data == nil {
		return
	}
	files := message.Data.([]interface{})
	//fmt.Println(files...)
	listfile := make([]model.FileInfo, 0, len(files))
	for _, f := range files {
		info := f.(map[string]interface{})
		mfile := model.FileInfo{}
		for k, v := range info {
			if k == "Name" {
				mfile.Name = v.(string)
			} else if k == "Path" {
				mfile.Path = v.(string)
			} else if k == "Size" {
				mfile.Size = v.(float64)
			} else if k == "Mode" {
				mfile.Mode = v.(string)
			} else if k == "ModTime" {
				mfile.ModTime = v.(string)
			} else if k == "IsDir" {
				mfile.IsDir = v.(bool)
			}
		}
		listfile = append(listfile, mfile)
	}
	if !l {
		List(listfile)
	} else {
		Table(listfile)
	}

}

func Table(list []model.FileInfo) {
	maxLength := 0
	for _, f := range list {
		if maxLength < len(f.Name) {
			maxLength = len(f.Name)
		}
		if maxLength < len(strconv.Itoa(int(f.Size))) {
			maxLength = len(strconv.Itoa(int(f.Size)))
		}
	}
	for _, f := range list {
		if len(strings.Trim(f.Name, " ")) > 0 {
			str := fmt.Sprintf("%s   %s %s   %s\n", f.Mode, handleStr(strconv.Itoa(int(f.Size)), maxLength), f.ModTime, handleStr(f.Name, maxLength))
			if f.IsDir {
				color.Bluep(str) //dir
			} else if strings.Contains(f.Mode, "x") {
				color.Greenp(str) //exe
			} else {
				color.White.Print(str) //file
			}
		}
	}
}

func handleStr(str string, length int) string {
	if len(strings.Trim(str, " ")) == 0 {
		return fillSpace(length)
	}
	if len(str) == length {
		return str
	}
	return str + fillSpace(length-len(str))
}

func fillSpace(length int) string {
	return strings.Join(make([]string, length), " ")
}

func List(list []model.FileInfo) {
	for _, f := range list {
		if len(strings.Trim(f.Name, " ")) > 0 {
			if f.IsDir {
				color.Bluep(" ", f.Name, "  ") //dir
			} else if strings.Contains(f.Mode, "x") {
				color.Greenp(" ", f.Name, "  ") //exe
			} else {
				color.White.Print(" ", f.Name, "  ") //file
			}
		}
	}
	fmt.Println("")
}
