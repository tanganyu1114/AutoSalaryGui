package mainws

import (
	"AutoSalaryGui/loginws"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
)

type MainGui struct {
	Window    *walk.MainWindow
	LoginPb   *walk.PushButton
	FilePb    *walk.PushButton
	ResetPb   *walk.PushButton
	SendPb    *walk.PushButton
	ShowLabel *walk.Label
	ShowView  *walk.WebView
}

func MainShow() {
	var mg MainGui
	var li loginws.LoginInfo
	ReadConf(&li)
	def := MainWindow{
		AssignTo: &mg.Window,
		Title:    "AutoSalary Gui  --ver0.1",
		MinSize:  Size{Width: 640, Height: 400},
		Size:     Size{Width: 640, Height: 400},
		Layout:   VBox{},
		Children: []Widget{
			PushButton{
				AssignTo: &mg.LoginPb,
				Text:     "登陆",
				Visible:  true,
				OnClicked: func() {
					//打开登陆配置界面
					if cmd, err := loginws.LoginWs(mg.Window, &li); err != nil {
						log.Print(err)
					} else if cmd == walk.DlgCmdOK {
						//保存用户信息
						SaveLogin(&li)
					}
				},
			},
			PushButton{
				AssignTo: &mg.FilePb,
				Text:     "请选择工资文件(xlsx)",
				Visible:  true,
				OnClicked: func() {
					//打开选择文件窗口，获取文件路径以及文件名
				},
			},
			PushButton{
				AssignTo: &mg.ResetPb,
				Text:     "重置",
				Visible:  true,
				OnClicked: func() {
					//重置选中的文件，清空预览窗口
				},
			},
			PushButton{
				AssignTo: &mg.SendPb,
				Text:     "发送邮件",
				Visible:  true,
				OnKeyUp: func(key walk.Key) {
					if key == walk.KeyS {
						mg.FilePb.SetText("")
					}
				},
			},
			Label{
				AssignTo: &mg.ShowLabel,
				Visible:  true,
				Text:     "预览邮件信息：",
			},
			WebView{
				AssignTo: &mg.ShowView,
				Visible:  true,
			},
		},
	}

	err := def.Create()
	if err != nil {
		fmt.Println(err)
	}
	_ = mg.Window.Run()
}
