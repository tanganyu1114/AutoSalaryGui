package mainws

import (
	"AutoSalaryGui/loginws"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type MainGui struct {
	Window    *walk.MainWindow
	LoginPb   *walk.PushButton
	FilePb    *walk.PushButton
	Fd        walk.FileDialog
	ResetPb   *walk.PushButton
	SendPb    *walk.PushButton
	LogPb     *walk.PushButton
	ShowLabel *walk.Label
	ShowView  *walk.WebView
}

var Mg MainGui

func MainShow() {

	//读取用户配置文件信息
	loginws.Li.ReadConf()

	Mg.Fd.Title = "请选择工资条文件"
	Mg.Fd.Filter = ".xlsx|*.xlsx"

	//返回dialog窗口
	def := MainWindow{
		AssignTo: &Mg.Window,
		Title:    "AutoSalary Gui  --ver0.1",
		MinSize:  Size{Width: 640, Height: 400},
		Size:     Size{Width: 640, Height: 400},
		Layout:   VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					PushButton{
						AssignTo:   &Mg.LoginPb,
						Text:       "登陆",
						Visible:    true,
						ColumnSpan: 2,
						OnClicked: func() {
							//打开登陆配置界面
							if cmd, err := loginws.LoginWs(Mg.Window); err != nil {
								WarnInfo(err.Error())
							} else if cmd == walk.DlgCmdOK {
								//保存用户信息
								loginws.Li.SaveLogin()
							}
						},
					},
					PushButton{
						AssignTo: &Mg.FilePb,
						Text:     "请选择工资文件(xlsx)",
						Visible:  true,
						OnClicked: func() {
							//打开选择文件窗口，获取文件路径以及文件名
							if cmd, err := Mg.Fd.ShowOpen(Mg.Window); err != nil {
								fmt.Println(err)
							} else if cmd {
								//成功获取路径并修改按钮名字
								Mg.FilePb.SetText(Mg.Fd.FilePath)
								Mg.FilePb.SetEnabled(false)
								//fmt.Println("filepath",Fd.FilePath,"title",Fd.Title)
								//确认选择后预览邮件发送效果
							}
						},
					},
					PushButton{
						AssignTo: &Mg.ResetPb,
						Text:     "重置",
						Visible:  true,
						OnClicked: func() {
							//重置选中的文件，清空预览窗口
							Mg.FilePb.SetEnabled(true)
							Mg.FilePb.SetText("请选择工资文件(xlsx)")
							Mg.Fd.FilePath = ""
							//fmt.Println("filepath",Fd.FilePath,"title",Fd.Title)
						},
					},
					PushButton{
						AssignTo: &Mg.SendPb,
						Text:     "发送邮件",
						Visible:  true,
						OnKeyUp: func(key walk.Key) {
							if key == walk.KeyS {
								Mg.FilePb.SetText("")
							}
						},
					},
					PushButton{
						AssignTo: &Mg.LogPb,
						Text:     "查看历史记录",
						Visible:  true,
						OnKeyUp: func(key walk.Key) {
							if key == walk.KeyS {
								Mg.FilePb.SetText("")
							}
						},
					},
					Label{
						AssignTo:   &Mg.ShowLabel,
						ColumnSpan: 2,
						Visible:    true,
						Text:       "预览邮件信息：",
					},
					WebView{
						AssignTo:   &Mg.ShowView,
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
	_ = Mg.Window.Run()
}
