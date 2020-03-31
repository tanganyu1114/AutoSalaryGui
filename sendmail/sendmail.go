package sendmail

import (
	"AutoSalaryGui/assets"
	"AutoSalaryGui/loginws"
	"AutoSalaryGui/setupmail"
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
	// xlsx表格内容
	Rows [][]string
	// 当前条数的邮件内容信息
	BodyInfo *bytes.Buffer
	// 收件人
	Touser string
	// 表头内容
	Header [][]string
	// 邮件条数
	Enum int
	// 合并的单元格信息
	Mcell    []excelize.MergeCell
	Si       *Sendinfo = &Sendinfo{Fmuser: loginws.Li.UserInfo}
	TempMail []byte
)

func init() {
	TempMail, _ = assets.Asset("source/mail_template.html")
}

// num 第几条表格信息
func GetMailInfo(num int) {
	var info string
	//rows := *xlsxptr
	//row := rows[num]
	num = num + len(Header)
	row := Rows[num]
	BodyInfo = new(bytes.Buffer)
	matched, _ := regexp.MatchString("\\w[-\\w.+]*@([A-Za-z0-9][-A-Za-z0-9]+\\.)+[A-Za-z]{2,14}", row[len(row)-1])
	if matched {
		Touser = row[len(row)-1]
	} else {
		info = "<meta http-equiv=\"Content-Type\" content=\"text/html; charset=utf-8\" /></br>获取邮箱失败,错误的邮箱格式,姓名：" + row[0] + " 邮箱信息:" + row[len(row)-1]
		BodyInfo = bytes.NewBufferString(info)
		return
	}
	//
	lnum := len(Header)
	rnum := len(Header[0])
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
				str += "<th style=\" text-align: center;font-size: 16px;height: 30px;width: 200px;border-bottom: 1px solid #999;border-right: 1px solid #999;\" rowspan=\"" + srspan + "\"colspan=\"" + scspan + "\">" + Header[j][i] + "</th>"
			} else {
				if Header[j][i] != "" {
					str += "<th style=\" text-align: center;font-size: 16px;height: 30px;width: 200px;border-bottom: 1px solid #999;border-right: 1px solid #999;\" rowspan=\"1\"  colspan=\"1\">" + Header[j][i] + "</th>"
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
		Spre:  setupmail.Mi.Prefix,
		Ssuf:  setupmail.Mi.Suffix,
		Ssign: setupmail.Mi.Sign,
		Stime: time.Now().Format("2006年01月02日"),
	}

	t := template.Must(template.New("mail_template.html").Parse(string(TempMail)))

	t.Execute(BodyInfo, Si)
}

// send the mail to user
// touser the mail receiver
// body  the mail information
func SendMail() (err error) {
	m := gomail.NewMessage()
	m.SetHeader("To", Touser)
	if setupmail.Mi.Alias != "" {
		m.SetAddressHeader("From", loginws.Li.UserInfo, setupmail.Mi.Alias)
	} else {
		m.SetHeader("From", loginws.Li.UserInfo)
	}
	m.SetHeader("Subject", setupmail.Mi.Title)
	m.SetBody("text/html", BodyInfo.String())

	d := gomail.NewDialer(loginws.Li.HostInfo, loginws.Li.PortInfo, loginws.Li.UserInfo, loginws.Li.PassInfo)

	// Send the email to user
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

// read the xlsx file
// return the file info / num / user
func ReadXlsx(xlsxPath string) (err error) {

	var index int
	var row []string
	xlsxFile, err := excelize.OpenFile(xlsxPath)
	if err != nil {
		return err
	}
	sheetname := xlsxFile.GetSheetName(sheetid)
	Rows, err = xlsxFile.GetRows(sheetname)
	if err != nil {
		return err
	}
	// 获取表头高度 index
	for index, row = range Rows {
		matched, _ := regexp.MatchString("\\w[-\\w.+]*@([A-Za-z0-9][-A-Za-z0-9]+\\.)+[A-Za-z]{2,14}", row[len(row)-1])
		if matched {
			break
		}
	}
	Mcell, _ = xlsxFile.GetMergeCells(sheetname)
	// table head
	Header = Rows[0:index]
	// info number
	Enum = len(Rows) - index
	//xlsptr = &rows
	return nil
}
