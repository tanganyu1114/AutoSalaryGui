package mainws

import (
	"github.com/lxn/walk"
)

func WarnInfo(str string) {
	walk.MsgBox(
		Mg.Window,
		"Error",
		str,
		walk.MsgBoxOK|walk.MsgBoxIconError)
}

func PromptInfo(str string) {
	walk.MsgBox(
		Mg.Window,
		"Info",
		str,
		walk.MsgBoxOK|walk.MsgBoxIconInformation)
}
