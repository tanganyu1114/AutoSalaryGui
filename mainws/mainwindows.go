package mainws

import (
	"AutoSalaryGui/filews"
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
	LogPb     *walk.PushButton
	ShowLabel *walk.Label
	ShowView  *walk.WebView
}

func MainShow() {
	var mg MainGui
	var li loginws.LoginInfo
	var xlsxpath string
	//读取用户配置文件信息
	ReadConf(mg, &li)

	//返回dialog窗口
	def := MainWindow{
		AssignTo: &mg.Window,
		Title:    "AutoSalary Gui  --ver0.1",
		MinSize:  Size{Width: 640, Height: 400},
		Size:     Size{Width: 640, Height: 400},
		Layout:   VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					PushButton{
						AssignTo:   &mg.LoginPb,
						Text:       "登陆",
						Visible:    true,
						ColumnSpan: 2,
						OnClicked: func() {
							//打开登陆配置界面
							if cmd, err := loginws.LoginWs(mg.Window, &li); err != nil {
								log.Print(err)
							} else if cmd == walk.DlgCmdOK {
								//保存用户信息
								SaveLogin(mg, &li)
							}
						},
					},
					PushButton{
						AssignTo: &mg.FilePb,
						Text:     "请选择工资文件(xlsx)",
						Visible:  true,
						OnClicked: func() {
							//打开选择文件窗口，获取文件路径以及文件名
							if cmd, err := filews.FileChoose(mg.Window, &xlsxpath); err != nil {
								log.Print(err)
							} else if cmd == walk.DlgCmdOK {
								//保存用户信息
								SaveLogin(mg, &li)
							}
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
					PushButton{
						AssignTo: &mg.LogPb,
						Text:     "查看历史记录",
						Visible:  true,
						OnKeyUp: func(key walk.Key) {
							if key == walk.KeyS {
								mg.FilePb.SetText("")
							}
						},
					},
					Label{
						AssignTo:   &mg.ShowLabel,
						ColumnSpan: 2,
						Visible:    true,
						Text:       "预览邮件信息：",
					},
					WebView{
						AssignTo:   &mg.ShowView,
						Visible:    true,
						ColumnSpan: 2,
					},
				},
			},
		},
	}

	err := def.Create()
	if err != nil {
		fmt.Println(err)
	}
	_ = mg.Window.Run()
}
