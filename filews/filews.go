package filews

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type XlsxInfo struct {
	filedlg  *walk.Dialog
	filePath *walk.DataBinder
	fileTree *walk.TreeView
	acceptPb *walk.PushButton
	cancelPb *walk.PushButton
}

//var XlsxFilePath	string

func FileChoose(wf walk.Form, xlsxpath *string) (int, error) {
	var fc XlsxInfo
	treeModel, err := NewDirectoryTreeModel()

	if err != nil {
		fmt.Println(err)
	}
	resfc := Dialog{
		AssignTo:      &fc.filedlg,
		Title:         "xlsx文件选择界面",
		DefaultButton: &fc.acceptPb,
		CancelButton:  &fc.cancelPb,
		DataBinder: DataBinder{
			AssignTo:       &fc.filePath,
			Name:           "xlsxpath",
			DataSource:     xlsxpath,
			ErrorPresenter: ToolTipErrorPresenter{},
		},
		Visible: true,
		MinSize: Size{Width: 400, Height: 500},
		Size:    Size{Width: 400, Height: 500},
		Layout:  VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						Text: "路径:",
					},
					LineEdit{
						CueBanner: "请输入xlsx文件路径",
						MinSize:   Size{Width: 300},
						Text:      Bind("xlsxpath"),
						OnTextChanged: func() {
							// Treeview根据路径跳转

						},
					},
					TreeView{
						AssignTo:   &fc.fileTree,
						ColumnSpan: 2,
						Model:      treeModel,
						MinSize:    Size{Width: 380, Height: 450},
						OnCurrentItemChanged: func() {
							dir := fc.fileTree.CurrentItem().(*Directory)
							xlsxpath = &dir.name
						},
					},
				},
			},
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					PushButton{
						AssignTo: &fc.acceptPb,
						Text:     "确认",
						//MinSize:Size{Width:100,Height:70},
						MaxSize: Size{Width: 150, Height: 100},
						//Alignment: AlignHNearVCenter,
						OnClicked: func() {
							//获取当前文件路径
							if err := fc.filePath.Submit(); err != nil {
								fmt.Println(err)
								return
							}
						},
					},
					PushButton{
						AssignTo: &fc.cancelPb,
						Text:     "取消",
						//MinSize:Size{Width:100,Height:70},
						MaxSize: Size{Width: 150, Height: 100},
						//Alignment: AlignHNearVCenter,
						OnClicked: func() {
							//关闭窗口
							fc.filedlg.Cancel()
						},
					},
				},
			},
		},
	}
	sint, err := resfc.Run(wf)
	return sint, err
}
