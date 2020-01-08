package sendmail

import (
	"AutoSalaryGui/loginws"
	"AutoSalaryGui/setmail"
	"bytes"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"gopkg.in/gomail.v2"
	"html/template"
	"regexp"
	"strconv"
	"time"
)

type Sendinfo struct {
	Title  string
	Fmuser string
	Touser string
	Sbody  template.HTML
	Spre   string
	Ssuf   string
	Ssign  string
	Stime  string
}

const (
	sheetid = 1
)

var (
	//	wd    *walk.MainWindow
	// 表头高度
	Head int = 2
	// 表头内容
	Entry [][]string
	// 邮件条数
	Enum int
	// 合并的单元格信息
	Mcell []excelize.MergeCell
	Si    *Sendinfo = &Sendinfo{Fmuser: loginws.Li.UserInfo}
)

func GetMailInfo(index int, rows [][]string) (buffer *bytes.Buffer) {
	var info string
	row := rows[Head+index]
	buffer = new(bytes.Buffer)
	matched, _ := regexp.MatchString("\\w[-\\w.+]*@([A-Za-z0-9][-A-Za-z0-9]+\\.)+[A-Za-z]{2,14}", row[len(row)-1])
	if matched {
		Si.Touser = row[len(row)-1]
	} else {

		info = "<meta http-equiv=\"Content-Type\" content=\"text/html; charset=utf-8\" /></br>获取邮箱失败,错误的邮箱格式,姓名：" + row[0] + " 邮箱信息:" + row[len(row)-1]
		buffer = bytes.NewBufferString(info)
		return
		//WarnInfo("获取邮箱失败,错误的邮箱格式,姓名："+row[0]+" 邮箱信息:"+row[len(row)-1])
	}
	//
	lnum := len(Entry)
	rnum := len(Entry[0])
	cnum := len(Mcell)
	mmcell := make(map[string]string, cnum)
	for _, m := range Mcell {
		mmcell[m.GetStartAxis()] = m.GetEndAxis()
	}

	for i := 0; i < rnum; i++ {
		//重置str为空
		str := ""
		for j := 0; j < lnum; j++ {

			var slnum, scnum, elnum, ecnum int //行号、列号

			//行列转换为cell坐标: A1/B1
			scell, _ := excelize.CoordinatesToCellName(i+1, j+1)

			if ecell, ok := mmcell[scell]; ok {
				slnum, scnum, _ = excelize.CellNameToCoordinates(scell)
				elnum, ecnum, _ = excelize.CellNameToCoordinates(ecell)
				srspan := strconv.Itoa(elnum - slnum + 1)
				scspan := strconv.Itoa(ecnum - scnum + 1)
				str += "<th style=\" text-align: center;font-size: 16px;height: 30px;width: 200px;border-bottom: 1px solid #999;border-right: 1px solid #999;\" rowspan=\"" + srspan + "\"colspan=\"" + scspan + "\">" + Entry[j][i] + "</th>"
			} else {
				if Entry[j][i] != "" {
					str += "<th style=\" text-align: center;font-size: 16px;height: 30px;width: 200px;border-bottom: 1px solid #999;border-right: 1px solid #999;\" rowspan=\"1\"  colspan=\"1\">" + Entry[j][i] + "</th>"
				}
			}
		}
		if i%2 == 0 {
			info += "<tr style=\" background-color: lightgrey;\">" + str + "<td style=\"text-align: center;font-size: 16px;height: 30px;width: 200px;border-bottom: 1px solid #999;border-right: 1px solid #999;\">" + row[i] + "</td>" + "</tr>"
		} else {
			info += "<tr style=\" background-color: white;\">" + str + "<td style=\"text-align: center;font-size: 16px;height: 30px;width: 200px;border-bottom: 1px solid #999;border-right: 1px solid #999;\">" + row[i] + "</td>" + "</tr>"
		}
	}

	body := template.HTML(info)

	Si = &Sendinfo{
		Sbody: body,
		Spre:  setmail.Mi.Prefix,
		Ssuf:  setmail.Mi.Suffix,
		Ssign: setmail.Mi.Sign,
		Stime: time.Now().Format("2006年01月02日"),
	}
	t := template.New("E:/Golang/src/AutoSalaryGui/sendmail/mail_template.html")
	t, _ = template.ParseFiles("E:/Golang/src/AutoSalaryGui/sendmail/mail_template.html")
	t.Execute(buffer, Si)
	return
}

/*func SendMail(rows [][]string) (view string) {


}*/

//发送所有
func SendAll(touser string, body string) (err error) {
	m := gomail.NewMessage()
	m.SetHeader("From", loginws.Li.UserInfo)
	m.SetHeader("To", touser)
	if setmail.Mi.Alias != "" {
		m.SetAddressHeader("Cc", loginws.Li.UserInfo, setmail.Mi.Alias)
	}
	m.SetHeader("Subject", setmail.Mi.Title)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(loginws.Li.HostInfo, loginws.Li.PortInfo, loginws.Li.UserInfo, loginws.Li.PassInfo)

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

//发送当前
/*func Send(xlsx string) (err error) {

}*/

func ReadXlsx(xlsxpath string) (err error, rows [][]string) {

	//改为Gui模式
	xlsFile, err := excelize.OpenFile(xlsxpath)
	if err != nil {
		return err, nil
	}
	sheetname := xlsFile.GetSheetName(sheetid)
	rows, err = xlsFile.GetRows(sheetname)
	if err != nil {
		return err, nil
	}
	Mcell, _ = xlsFile.GetMergeCells(sheetname)
	Entry = rows[0:Head]
	Enum = len(rows) - Head
	return nil, rows
}
