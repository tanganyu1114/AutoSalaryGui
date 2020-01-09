package mainws

import (
	"AutoSalaryGui/sendmail"
	"github.com/lxn/walk"
	"io/ioutil"
	"path/filepath"
	"strconv"
)

func SetLinelabel() {
	Mg.NumInfo.SetText("第" + strconv.Itoa(index+1) + "条/总共" + strconv.Itoa(num) + "条")
}

func View() {
	//bind webview info
	if rows != nil {
		touser, mailinfo = sendmail.GetMailInfo(index, rows)
		err := ioutil.WriteFile(TempFile, mailinfo.Bytes(), 0660)
		if err != nil {
			WarnInfo(err.Error())
		}
		tmppath, _ := filepath.Abs(TempFile)
		//fmt.Println(tmppath)
		Mg.ShowView.SetURL(tmppath)
	} else {
		WarnInfo("请选择工资条文件(.xlsx)")
	}
}

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
