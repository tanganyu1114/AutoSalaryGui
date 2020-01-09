package mainws

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type PromptSend struct {
	pmpDlg  *walk.Dialog
	pmpShow *walk.Label
	pmpPb   *walk.PushButton
}

var ps PromptSend

func PromptSendInfo(wf walk.Form) (int, error) {

	pmptDlg := Dialog{
		AssignTo:      &ps.pmpDlg,
		Title:         "Info",
		Icon:          "source/logo.png",
		DefaultButton: &ps.pmpPb,
		CancelButton:  &ps.pmpPb,
		Visible:       true,
		MinSize:       Size{Width: 400, Height: 320},
		Size:          Size{Width: 400, Height: 320},
		Layout:        VBox{},
		Children: []Widget{
			Label{
				Text: "邮件发送中，请等待！",
			},
			Label{
				AssignTo: &ps.pmpShow,
			},
			PushButton{
				Text:      "Ok",
				AssignTo:  &ps.pmpPb,
				OnClicked: func() { ps.pmpDlg.Cancel() },
			},
		},
	}
	sint, err := pmptDlg.Run(wf)
	return sint, err
}
