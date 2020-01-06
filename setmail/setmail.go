package setmail

import (
	"AutoSalaryGui/loginws"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
)

type SetGui struct {
	setmailDlg *walk.Dialog
	acceptPb   *walk.PushButton
	cancelPb   *walk.PushButton
	loginDb    *walk.DataBinder
}

type MailStruct struct {
	Title     string
	Aliasuser string
	Prefix    string
	Suffix    string
	Sign      string
}

var (
	sg SetGui
	ms MailStruct
)

func SetMail(wf walk.Form) (int, error) {

	resmail := Dialog{
		AssignTo:      &sg.setmailDlg,
		Title:         "邮件配置界面",
		DefaultButton: &sg.acceptPb,
		CancelButton:  &sg.cancelPb,
		DataBinder: DataBinder{
			AssignTo:       &sg.loginDb,
			Name:           "LoginInfo",
			DataSource:     ms,
			ErrorPresenter: ToolTipErrorPresenter{},
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
						Text:      Bind("Aliasuser"),
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
							if err := sg.loginDb.Submit(); err != nil {
								log.Print(err)
								return
							}

							sg.setmailDlg.Accept()
						},
					},
					PushButton{
						AssignTo:  &sg.cancelPb,
						Text:      "Cancel",
						OnClicked: func() { sg.setmailDlg.Cancel() },
					},
				},
			},
		},
	}
	sint, err := resmail.Run(wf)
	return sint, err
}
