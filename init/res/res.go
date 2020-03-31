package res

// 全局变量

import (
	"AutoSalaryGui/assets"
	"io/ioutil"
	"os"
)

var TEMP string = os.Getenv("TEMP")
var Logo, Viewbg []byte

func init() {
	Logo, _ = assets.Asset("source/logo.png")
	Viewbg, _ = assets.Asset("source/view.jpg")
	if !FileExist(TEMP + "\\logo.png") {
		ioutil.WriteFile(TEMP+"\\logo.png", Logo, 0660)
	}
	if !FileExist(TEMP + "\\view.jpg") {
		ioutil.WriteFile(TEMP+"\\view.jpg", Viewbg, 0660)
	}
}

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}
