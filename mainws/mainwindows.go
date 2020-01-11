package mainws

import (
	"AutoSalaryGui/loginws"
	"AutoSalaryGui/sendmail"
	"AutoSalaryGui/setmail"
	"AutoSalaryGui/source"
	"bytes"
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

const (
	TempFile = "temp_mail_file.html"
)

var (
	Mg MainGui
	// mail information ptr
	bodyptr *bytes.Buffer
	// mail number
	index int
	// 视图背景
	viewpath     string
	Logo, viewbg []byte
)

func init() {
	Logo, _ = source.Asset("source/logo.png")
	viewbg, _ = source.Asset("source/view.jpg")
	if !loginws.FileExist("logo.png") {
		ioutil.WriteFile("logo.png", Logo, 0660)
	}
	if !loginws.FileExist("view.jpg") {
		ioutil.WriteFile("view.jpg", viewbg, 0660)
	}
}

func MainShow() {

	//返回dialog窗口
	def := MainWindow{
		AssignTo: &Mg.Window,
		Icon:     "logo.png",
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
								Mg.LoginPb.SetText(loginws.Li.UserInfo + "(已配置)")
							}
						},
					},
					PushButton{
						AssignTo: &Mg.SendAllPb,
						Text:     "发送全部工资条",
						Visible:  true,
						Enabled:  false,
						OnClicked: func() {
							//点击发送邮件
							//sendmail.SendAll(rows,mailinfo)
							if sendmail.Rows != nil {
								PromptSendInfo(Mg.Window)
								Mg.SendAllPb.SetEnabled(false)
								var res string
								for index = 0; index < sendmail.Enum; index++ {
									sendmail.GetMailInfo(index)
									err := sendmail.SendMail()
									if err != nil {
										res = err.Error()
									}
									ps.pmpShow.SetText("正在发送第" + strconv.Itoa(index+1) + "条/总共" + strconv.Itoa(sendmail.Enum) + "条")
								}
								ps.pmpPb.Clicked()
								if res != "" {
									WarnInfo(res)
								} else {
									PromptInfo("邮件发送完成!")
								}
							} else {
								WarnInfo("请选择工资文件（xlsx）")
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
								// read the xlsx file
								err = sendmail.ReadXlsx(Mg.Fd.FilePath)
								if err != nil {
									WarnInfo(err.Error())
								}
								// set the numbedit range
								Mg.NumEdit.SetRange(1, float64(sendmail.Enum))
								if sendmail.Enum > 1 {
									Mg.NextPd.SetEnabled(true)
								}
								// show view
								SetLinelabel()
								View()
								// set enable the pb
								Mg.SendSiglePb.SetEnabled(true)
								Mg.SendAllPb.SetEnabled(true)
								Mg.GotoPb.SetEnabled(true)
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
							if sendmail.Rows != nil {
								sendmail.Rows = make([][]string, 0)
							}
							Mg.ShowView.SetURL(viewpath)
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
						Text:     "查看日志",
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
						Text:      "第0条/总共0条",
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
							if index > sendmail.Enum-1 {
								index = sendmail.Enum - 1
							}
							if index == sendmail.Enum-1 {
								Mg.NextPd.SetEnabled(false)
							}
							SetLinelabel()
							View()
						},
					},
					NumberEdit{
						AssignTo: &Mg.NumEdit,
						Decimals: 0,
					},
					PushButton{
						AssignTo: &Mg.GotoPb,
						Text:     "跳转",
						Visible:  true,
						Enabled:  false,
						OnClicked: func() {
							//index +1 , if index=Enum disable
							index = int(Mg.NumEdit.Value()) - 1
							if index < 0 {
								index = 0
							}
							if index == 0 {
								Mg.ForwardPb.SetEnabled(false)
							} else {
								Mg.ForwardPb.SetEnabled(true)
							}
							if index > sendmail.Enum-1 {
								index = sendmail.Enum - 1
							}
							if index == sendmail.Enum-1 {
								Mg.NextPd.SetEnabled(false)
							} else {
								Mg.NextPd.SetEnabled(true)
							}
							SetLinelabel()
							View()
						},
					},
					PushButton{
						AssignTo: &Mg.SendSiglePb,
						Text:     "发送当前邮件",
						Visible:  true,
						Enabled:  false,
						OnClicked: func() {
							// send the mail info
							if sendmail.Rows != nil {
								err := sendmail.SendMail()
								if err != nil {
									WarnInfo(err.Error())
								} else {
									PromptInfo("邮件发送成功！")
								}
							} else {
								WarnInfo("请选择工资文件（xlsx）")
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
	//初始化窗口，加载数据信息
	initws()
	_ = Mg.Window.Run()
	//loadInfo()

}

//初始化窗口，加载数据信息
func initws() {
	//读取用户配置文件信息
	if err := loginws.Li.ReadConf(); err != nil {
		WarnInfo(err.Error())
	} else {
		if loginws.Li.UserInfo != "" {
			Mg.LoginPb.SetText(loginws.Li.UserInfo + "(已配置)")
		}
	}
	//读取邮件配置信息
	if err := setmail.Mi.ReadMailConf(); err != nil {
		WarnInfo(err.Error())
	}

	Mg.Fd.Title = "请选择工资条文件"
	Mg.Fd.Filter = "*.xlsx|*.xlsx"
	viewpath, _ = os.Getwd()
	viewpath = viewpath + "/view.jpg"
	Mg.ShowView.SetURL(viewpath)
}

// 读取xlsx文件，加载读取后的界面信息
/*func loadxlsx(){

}*/
