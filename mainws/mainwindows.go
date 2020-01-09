package mainws

import (
	"AutoSalaryGui/loginws"
	"AutoSalaryGui/sendmail"
	"AutoSalaryGui/setmail"
	"bytes"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
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

// xlsx内容
var rows [][]string

// 表头高度
var head int = 2
var touser string
var mailinfo *bytes.Buffer
var index int = 0
var num int = 0

func init() {
	//读取用户配置文件信息
	if err := loginws.Li.ReadConf(); err != nil {
		WarnInfo(err.Error())
	}
	//读取邮件配置信息
	if err := setmail.Mi.ReadMailConf(); err != nil {
		WarnInfo(err.Error())
	}
}

func MainShow() {

	Mg.Fd.Title = "请选择工资条文件"
	Mg.Fd.Filter = "*.xlsx|*.xlsx"
	viewbg, _ := filepath.Abs("source/view.jpg")
	num = len(rows) - head
	if num < 0 {
		num = 1
	}

	//返回dialog窗口
	def := MainWindow{
		AssignTo: &Mg.Window,
		Icon:     "source/logo.png",
		//Background: BitmapBrush{Image:"source/view.jpg",},
		Title:   "AutoSalary Gui  --ver0.1",
		MinSize: Size{Width: 640, Height: 720},
		Size:    Size{Width: 640, Height: 720},
		Layout:  VBox{},
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
							//sendmail.SendAll(rows,mailinfo)
							PromptSendInfo(Mg.Window)
							Mg.SendAllPb.SetEnabled(false)
							var res string
							for index = 0; index < len(rows)-head; index++ {
								touser, mailinfo = sendmail.GetMailInfo(index, rows)
								err := sendmail.SendSigle(touser, mailinfo)
								if err != nil {
									res = err.Error()
								}
								ps.pmpShow.SetText("正在发送第" + strconv.Itoa(index+1) + "条/总共" + strconv.Itoa(num) + "条")
							}
							ps.pmpPb.Clicked()
							if res != "" {
								WarnInfo(res)
							} else {
								PromptInfo("邮件发送完成!")
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
								err, rows = sendmail.ReadXlsx(Mg.Fd.FilePath)
								if err != nil {
									WarnInfo(err.Error())
								}
								num = len(rows) - head
								if num > 1 {
									Mg.NextPd.SetEnabled(true)
								}
								SetLinelabel()
								View()
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
							Mg.ShowView.SetURL(viewbg)
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
								if err := setmail.Mi.SaveMailConf(); err != nil {
									WarnInfo(err.Error())
								}
								View()
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
					WebView{
						AssignTo:   &Mg.ShowView,
						Visible:    true,
						URL:        viewbg,
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
						Enabled:  false,
						OnClicked: func() {
							//index -1 , if index=0 disable
							index -= index
							/*							switch {
														case index<0:
															index=0
															fallthrough
														case index==0:
															Mg.ForwardPb.SetVisible(false)
														}*/
							Mg.NextPd.SetEnabled(true)
							if index < 0 {
								index = 0
							}
							if index == 0 {
								Mg.ForwardPb.SetEnabled(false)
							}
							SetLinelabel()
							View()
						},
					},
					Label{
						AssignTo:  &Mg.NumInfo,
						Alignment: AlignHCenterVCenter,
						Visible:   true,
						Text:      "第" + strconv.Itoa(index+1) + "条/总共" + strconv.Itoa(num) + "条",
					},
					PushButton{
						AssignTo: &Mg.NextPd,
						Text:     "下一条",
						Visible:  true,
						Enabled:  false,
						OnClicked: func() {
							//index +1 , if index=Enum disable
							index += 1
							Mg.ForwardPb.SetEnabled(true)
							if index > len(rows)-head-1 {
								index = len(rows) - head - 1
							}
							if index == len(rows)-head-1 {
								Mg.NextPd.SetEnabled(false)
							}
							SetLinelabel()
							View()
						},
					},
					NumberEdit{
						AssignTo: &Mg.NumEdit,
						MinValue: 1,
						MaxValue: float64(num),
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
							// send the mail info
							err := sendmail.SendSigle(touser, mailinfo)
							if err != nil {
								WarnInfo(err.Error())
							} else {
								PromptInfo("邮件发送成功！")
							}
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
