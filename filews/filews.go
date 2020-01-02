package filews

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type XlsxInfo struct {
	fileDlg  *walk.Dialog
	filePath *walk.DataBinder
	fileTree *walk.TreeView
	fileView *walk.TableView
	acceptPb *walk.PushButton
	cancelPb *walk.PushButton
	spLitter *walk.Splitter
}

//var XlsxFilePath	string

func FileChoose(wf walk.Form, xlsxpath *string) (int, error) {
	var fc XlsxInfo
	treeModel, err := NewDirectoryTreeModel()
	xlsxModel := NewFileInfoModel()
	if err != nil {
		fmt.Println(err)
	}
	resfc := Dialog{
		AssignTo:      &fc.fileDlg,
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
		MinSize: Size{Width: 500, Height: 500},
		Size:    Size{Width: 500, Height: 500},
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
					HSplitter{
						AssignTo:   &fc.spLitter,
						ColumnSpan: 2,
						Children: []Widget{
							TreeView{
								AssignTo:   &fc.fileTree,
								ColumnSpan: 1,
								Model:      treeModel,
								MinSize:    Size{Width: 165, Height: 450},
								OnCurrentItemChanged: func() {
									//dir := fc.fileTree.CurrentItem().(*Directory)
									//xlsxpath = &dir.name
									dir := fc.fileTree.CurrentItem().(*Directory)
									if err := xlsxModel.SetDirPath(dir.Path()); err != nil {
										//mainws.WarnInfo(wf,err.Error())
										fmt.Println(err)
									}
								},
							},
							TableView{
								AssignTo: &fc.fileView,
								//StretchFactor: 2,
								Columns: []TableViewColumn{
									TableViewColumn{
										DataMember: "Name",
										Width:      165,
									},
									TableViewColumn{
										DataMember: "Size",
										Format:     "%d",
										Alignment:  AlignFar,
										Width:      60,
									},
									TableViewColumn{
										DataMember: "Modified",
										Format:     "2006-01-02 15:04:05",
										Width:      120,
									},
								},
								Model: xlsxModel,
								OnCurrentIndexChanged: func() {
									/*									var url string
																		if index := fc.fileView.CurrentIndex(); index > -1 {
																			name := xlsxModel.items[index].Name
																			dir := fc.fileTree.CurrentItem().(*Directory)
																			url = filepath.Join(dir.Path(), name)
																		}
									*/
									//webView.SetURL(url)
								},
							},
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
							fc.fileDlg.Cancel()
						},
					},
				},
			},
		},
	}
	sint, err := resfc.Run(wf)
	return sint, err
}
