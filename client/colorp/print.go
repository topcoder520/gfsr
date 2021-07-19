package colorp

import (
	"github.com/gookit/color"
	"github.com/topcoder520/gfsr/client/config"
)

func BluePrint(msg string) {
	color.Bluep(msg)
}

func BluePrintln(msg string) {
	color.Blueln(msg)
}

func GreenPrint(msg string) {
	color.Greenp(msg)
}

func GreenPrintln(msg string) {
	color.Greenln(msg)
}

func WhitePrint(msg string) {
	color.White.Print(msg)
}

func WhitePrintln(msg string) {
	color.White.Println(msg)
}

func ErrorPrintln(err string) {
	color.Redln(err)
}

func LineHeaderInfo() {
	color.LightGreen.Printf(config.UserName)
	color.White.Printf(":")
	p := "~" + config.CurrentPath
	if p == "~/" {
		p = "~"
	}
	color.Blue.Printf(p)
	color.White.Printf("$ ")
}
