package setmail

import (
	"AutoSalaryGui/loginws"
	"encoding/json"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type SetGui struct {
	setDlg   *walk.Dialog
	acceptPb *walk.PushButton
	cancelPb *walk.PushButton
	setDb    *walk.DataBinder
}

type MailInfo struct {
	Title  string
	Alias  string
	Prefix string
	Suffix string
	Sign   string
}

type SetMailConf interface {
	SaveMailConf()
	ReadMailConf()
}

const MailConf = "mail.config"

var (
	sg SetGui
	Mi = new(MailInfo)
)

func SetMail(wf walk.Form) (int, error) {

	resmail := Dialog{
		AssignTo:      &sg.setDlg,
		Title:         "邮件配置界面",
		Icon:          "logo.png",
		DefaultButton: &sg.acceptPb,
		CancelButton:  &sg.cancelPb,
		DataBinder: DataBinder{
			AssignTo:        &sg.setDb,
			Name:            "MailInfo",
			DataSource:      Mi,
			AutoSubmit:      true,
			AutoSubmitDelay: time.Second * 3,
			ErrorPresenter:  ToolTipErrorPresenter{},
		},
		Visible: true,
		MinSize: Size{Width: 500, Height: 640},
		Size:    Size{Width: 500, Height: 640},
		Layout:  VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						Text: "邮件主题:",
					},
					LineEdit{
						CueBanner: "请输入邮件主题",
						Text:      Bind("Title"),
					},
					Label{
						Text: "发件人:",
					},
					LineEdit{
						CueBanner: "default: " + loginws.Li.UserInfo,
						Text:      Bind("Alias"),
					},
					Label{
						ColumnSpan: 2,
						Text:       "邮件内容(表格前):",
					},
					TextEdit{
						ColumnSpan: 2,
						VScroll:    true,
						MinSize:    Size{Height: 200, Width: 480},
						Text:       Bind("Prefix"),
					},
					Label{
						ColumnSpan: 2,
						Text:       "邮件内容(表格后):",
					},
					TextEdit{
						ColumnSpan: 2,
						VScroll:    true,
						MinSize:    Size{Height: 200, Width: 480},
						Text:       Bind("Suffix"),
					},
					Label{
						Text: "签名:",
					},
					LineEdit{
						CueBanner: "邮件落款签名",
						Text:      Bind("Sign"),
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo: &sg.acceptPb,
						Text:     "OK",
						OnClicked: func() {
							if err := sg.setDb.Submit(); err != nil {
								log.Print(err)
								return
							}

							sg.setDlg.Accept()
						},
					},
					PushButton{
						AssignTo:  &sg.cancelPb,
						Text:      "Cancel",
						OnClicked: func() { sg.setDlg.Cancel() },
					},
				},
			},
		},
	}
	sint, err := resmail.Run(wf)
	return sint, err
}

func (mi *MailInfo) SaveMailConf() (err error) {
	//	fmt.Println(mi)
	data, _ := json.Marshal(mi)
	err = ioutil.WriteFile(MailConf, data, 0660)
	if err != nil {
		return
	}
	return nil
}

func (mi *MailInfo) ReadMailConf() (err error) {
	path, _ := os.Getwd()
	filepath := path + "/" + MailConf
	if FileExist(filepath) {
		filePtr, err := os.Open(MailConf)
		if err != nil {
			//fmt.Println("读取用户信息失败！")
			return err
		}
		defer filePtr.Close()
		readInfo := json.NewDecoder(filePtr)
		err = readInfo.Decode(mi)
		if err != nil {
			//fmt.Println("存储用户信息格式错误！")
			return err
		}
	}
	return nil
}

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}
