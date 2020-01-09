package loginws

import (
	"encoding/base64"
	"encoding/json"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"io/ioutil"
	"log"
	"os"
)

type LoginGui struct {
	loginDlg *walk.Dialog
	acceptPb *walk.PushButton
	cancelPb *walk.PushButton
	loginDb  *walk.DataBinder
}

type Login interface {
	SaveLogin()
	ReadConf()
}

type LoginInfo struct {
	UserInfo  string
	PassInfo  string
	HostInfo  string
	PortInfo  int
	ValidInfo bool
}

const LoginConf = "login.config"

var (
	Li *LoginInfo = &LoginInfo{PortInfo: 465}
	lg LoginGui
)

func LoginWs(wf walk.Form) (int, error) {

	reslg := Dialog{
		AssignTo:      &lg.loginDlg,
		Title:         "Email邮箱登陆界面",
		Icon:          "source/logo.png",
		DefaultButton: &lg.acceptPb,
		CancelButton:  &lg.cancelPb,
		DataBinder: DataBinder{
			AssignTo:       &lg.loginDb,
			Name:           "LoginInfo",
			DataSource:     Li,
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
						CueBanner:    "请输入您的密码或授权码",
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
						AssignTo: &lg.acceptPb,
						Text:     "OK",
						OnClicked: func() {
							if err := lg.loginDb.Submit(); err != nil {
								log.Print(err)
								return
							}

							lg.loginDlg.Accept()
						},
					},
					PushButton{
						AssignTo:  &lg.cancelPb,
						Text:      "Cancel",
						OnClicked: func() { lg.loginDlg.Cancel() },
					},
				},
			},
		},
	}
	sint, err := reslg.Run(wf)
	return sint, err
}

func (li *LoginInfo) SaveLogin() (err error) {

	//	var filePtr *os.File
	srcpass := li.PassInfo
	passbyte := []byte(li.PassInfo)
	//加密
	encoded := base64.StdEncoding.EncodeToString(passbyte)
	//解密
	//decoded, _ := base64.StdEncoding.DecodeString(encoded)
	li.PassInfo = encoded
	data, _ := json.Marshal(li)

	err = ioutil.WriteFile(LoginConf, data, 0660)
	if err != nil {
		return
	}
	li.PassInfo = srcpass
	return nil
}

func (li *LoginInfo) ReadConf() (err error) {
	path, _ := os.Getwd()
	filepath := path + "/" + LoginConf
	if FileExist(filepath) {
		filePtr, err := os.Open(LoginConf)
		if err != nil {
			//fmt.Println("读取用户信息失败！")
			return err
		}
		defer filePtr.Close()
		readInfo := json.NewDecoder(filePtr)
		err = readInfo.Decode(li)
		if err != nil {
			//fmt.Println("存储用户信息格式错误！")
			return err
		}
		decoded, _ := base64.StdEncoding.DecodeString(li.PassInfo)
		li.PassInfo = string(decoded)
	}
	return nil
}

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

//登陆校验
/*func (lg *LoginInfo)LoginValid()(valid bool,err error){



}*/
