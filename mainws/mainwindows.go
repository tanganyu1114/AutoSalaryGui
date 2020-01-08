package mainws

import (
	"AutoSalaryGui/loginws"
	"AutoSalaryGui/sendmail"
	"AutoSalaryGui/setmail"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

type MainGui struct {
	Window      *walk.MainWindow
	LoginPb     *walk.PushButton
	FilePb      *walk.PushButton
	SetmailPb   *walk.PushButton
	ResetPb     *walk.PushButton
	SendAllPb   *walk.PushButton
	LogPb       *walk.PushButton
	ShowLabel   *walk.Label
	ShowView    *walk.WebView
	ViewPb      *walk.PushButton
	NumInfo     *walk.Label
	ForwardPb   *walk.PushButton
	NextPd      *walk.PushButton
	NumEdit     *walk.NumberEdit
	GotoPb      *walk.PushButton
	SendSiglePb *walk.PushButton
	Fd          walk.FileDialog
}

const TempFile = "temp_mail_file.html"

var Mg MainGui
var rows [][]string
var index int = 0

func init() {
	//读取用户配置文件信息
	loginws.Li.ReadConf()
	//读取邮件配置信息
	setmail.Mi.ReadMailConf()
}

func MainShow() {

	Mg.Fd.Title = "请选择工资条文件"
	Mg.Fd.Filter = ".xlsx|*.xlsx"

	//返回dialog窗口
	def := MainWindow{
		AssignTo: &Mg.Window,
		Title:    "AutoSalary Gui  --ver0.1",
		MinSize:  Size{Width: 640, Height: 720},
		Size:     Size{Width: 640, Height: 720},
		Layout:   VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					PushButton{
						AssignTo:   &Mg.LoginPb,
						Text:       "登陆",
						Visible:    true,
						ColumnSpan: 1,
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
						AssignTo: &Mg.SendAllPb,
						Text:     "发送全部工资条",
						Visible:  true,
						OnClicked: func() {
							//点击发送邮件
							//sendmail.SendMail(rows)
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
								err, rows = sendmail.ReadXlsx(Mg.Fd.FilePath)
								if err != nil {
									WarnInfo(err.Error())
								}
								//fmt.Println(rows)
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
						AssignTo: &Mg.SetmailPb,
						Text:     "邮件配置",
						Visible:  true,
						OnClicked: func() {

							if cmd, err := setmail.SetMail(Mg.Window); err != nil {
								WarnInfo(err.Error())
							} else if cmd == walk.DlgCmdOK {
								//保存邮件配置信息
								setmail.Mi.SaveMailConf()
							}
						},
					},

					PushButton{
						AssignTo: &Mg.LogPb,
						Text:     "查看历史记录",
						Visible:  true,
						OnClicked: func() {
							dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
							fmt.Println(dir)
						},
					},
					Label{
						AssignTo:   &Mg.ShowLabel,
						ColumnSpan: 1,
						Visible:    true,
						Text:       "预览邮件信息：",
					},
					PushButton{
						AssignTo:  &Mg.ViewPb,
						Text:      "预览",
						Visible:   true,
						Alignment: AlignHFarVNear,
						OnClicked: func() {
							//bind webview info
							if rows != nil {
								mailinfo := sendmail.GetMailInfo(index, rows)
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
						},
					},
					WebView{
						AssignTo:   &Mg.ShowView,
						Visible:    true,
						ColumnSpan: 2,
					},
				},
			},
			Composite{
				Layout: Grid{Columns: 3},
				Children: []Widget{
					PushButton{
						AssignTo: &Mg.ForwardPb,
						Text:     "上一条",
						Visible:  true,
						OnClicked: func() {
							//index -1 , if index=0 disable
						},
					},
					Label{
						AssignTo:  &Mg.NumInfo,
						Alignment: AlignHCenterVCenter,
						Visible:   true,
						Text:      "第" + strconv.Itoa(index) + "条/总共" + strconv.Itoa(sendmail.Enum) + "条",
					},
					PushButton{
						AssignTo: &Mg.NextPd,
						Text:     "下一条",
						Visible:  true,
						OnClicked: func() {
							//index +1 , if index=Enum disable

						},
					},
					NumberEdit{
						AssignTo: &Mg.NumEdit,
						Decimals: 0,
						Value:    index,
					},
					PushButton{
						AssignTo: &Mg.GotoPb,
						Text:     "跳转",
						Visible:  true,
						OnClicked: func() {
							//index +1 , if index=Enum disable

						},
					},
					PushButton{
						AssignTo: &Mg.SendSiglePb,
						Text:     "发送当前邮件",
						Visible:  true,
						OnClicked: func() {
							//index +1 , if index=Enum disable

						},
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
