package loginws

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
)

type LoginGui struct {
	Dlg       *walk.Dialog
	UserLable *walk.Label
	PassLable *walk.Label
	UserEdit  *walk.LineEdit
	PassEdit  *walk.LineEdit
}

type LoginInfo struct {
	UserInfo string
	PassInfo string
	HostInfo string
	PortInfo int
}

func LoginWs(wf walk.Form, li *LoginInfo) (int, error) {
	var lg LoginGui
	var db *walk.DataBinder
	var acceptPB, cancelPB *walk.PushButton

	reslg := Dialog{
		AssignTo:      &lg.Dlg,
		Title:         "Email邮箱登陆界面",
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		DataBinder: DataBinder{
			AssignTo:       &db,
			Name:           "LoginInfo",
			DataSource:     li,
			ErrorPresenter: ToolTipErrorPresenter{},
		},
		Visible: true,
		MinSize: Size{Width: 400, Height: 320},
		Size:    Size{Width: 400, Height: 320},
		Layout:  VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						Text: "账号:",
					},
					LineEdit{
						CueBanner: "请输入您的邮箱账号",
						Text:      Bind("UserInfo"),
					},
					Label{
						Text: "密码:",
					},
					LineEdit{
						PasswordMode: true,
						CueBanner:    "请输入您的邮箱密码",
						Text:         Bind("PassInfo"),
					},
					Label{
						Text: "服务器地址:",
					},
					LineEdit{
						CueBanner: "请输入邮箱服务器地址",
						Text:      Bind("HostInfo"),
					},
					Label{
						Text: "端口号:",
					},
					NumberEdit{
						Value:    Bind("PortInfo", Range{0, 9999}),
						Decimals: 0,
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo: &acceptPB,
						Text:     "OK",
						OnClicked: func() {
							if err := db.Submit(); err != nil {
								log.Print(err)
								return
							}

							lg.Dlg.Accept()
						},
					},
					PushButton{
						AssignTo:  &cancelPB,
						Text:      "Cancel",
						OnClicked: func() { lg.Dlg.Cancel() },
					},
				},
			},
		},
	}
	sint, err := reslg.Run(wf)
	return sint, err
}
