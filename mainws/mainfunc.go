package mainws

import (
	"github.com/lxn/walk"
)

func ChooseFile() {

}

func Reset() {

}

func WarnInfo(str string) {
	walk.MsgBox(
		Mg.Window,
		"Error",
		str,
		walk.MsgBoxOK|walk.MsgBoxIconError)
}
