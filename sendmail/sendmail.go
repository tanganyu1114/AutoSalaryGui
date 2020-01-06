package sendmail

import (
	"AutoSalaryGui/mainws"
	"bytes"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/lxn/walk"
	"gopkg.in/gomail.v2"
	"html/template"
	"os"
	"regexp"
	"strconv"
	"time"
)

const (
	sheetid = 0
)

var (
	warn  *walk.MainWindow
	Entry [][]string
	Enum  int
	Nmail int = 0
	Head  int = 1
	Mcell []excelize.MergeCell
)

type Sendinfo struct {
	Title  string
	Touser string

	Sbody template.HTML
	Spre  string
	Ssuf  string
	Stime string
}

func sendMail(rows [][]string) {
	var touser string

	for _, row := range rows[Head:] {
		buffer := new(bytes.Buffer)
		matched, _ := regexp.MatchString("\\w[-\\w.+]*@([A-Za-z0-9][-A-Za-z0-9]+\\.)+[A-Za-z]{2,14}", row[len(row)-1])
		if matched {
			touser = row[len(row)-1]

		} else {
			fmt.Println("获取邮箱失败,错误的邮箱格式,姓名：", row[0], "邮箱信息", row[len(row)-1])
			continue
		}
		body := template.HTML(Emailbody(Entry, row))

		Sendinfo := Sendinfo{
			Sbody: body,
			Spre:  Emailconfig.Mailinfo.Pre,
			Ssuf:  Emailconfig.Mailinfo.Suf,
			Stime: time.Now().Format("2006年01月02日"),
		}

		t := template.New("email_template.html")
		t, _ = template.ParseFiles("email_template.html")
		t.Execute(buffer, Sendinfo)
		Nmail += 1
		if Issend == 1 {
			fmt.Println("总共[", Enum, "]封邮件,开始发送第[", Nmail, "]封邮件,发送给", touser)
			send(touser, buffer.String())
		}
		if Issend == 2 {
			fmt.Println("开始打印预览文件")
			cfile, err := os.Create("test.html")
			defer cfile.Close()
			if err != nil {
				fmt.Println(err.Error())
			} else {
				_, err = cfile.Write([]byte(buffer.String()))
				checkErr(err)
			}
			fmt.Println("文件打印完毕：test.html")
			fmt.Println("输入回车字符退出程序 ！")
		}
	}
	fmt.Println("总共[", Enum, "]封邮件,已全部发送完毕 !")
	fmt.Printf("\t")
	fmt.Println("输入回车字符退出程序 ！")

}

func Emailbody(entry [][]string, row []string) (info string) {
	lnum := len(entry)
	rnum := len(entry[0])
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
				str += "<th style=\" text-align: center;font-size: 16px;height: 30px;width: 200px;border-bottom: 1px solid #999;border-right: 1px solid #999;\" rowspan=\"" + srspan + "\"colspan=\"" + scspan + "\">" + entry[j][i] + "</th>"
			} else {
				if entry[j][i] != "" {
					str += "<th style=\" text-align: center;font-size: 16px;height: 30px;width: 200px;border-bottom: 1px solid #999;border-right: 1px solid #999;\" rowspan=\"1\"  colspan=\"1\">" + entry[j][i] + "</th>"
				}
			}
		}
		if i%2 == 0 {
			info += "<tr style=\" background-color: lightgrey;\">" + str + "<td style=\"text-align: center;font-size: 16px;height: 30px;width: 200px;border-bottom: 1px solid #999;border-right: 1px solid #999;\">" + row[i] + "</td>" + "</tr>"
		} else {
			info += "<tr style=\" background-color: white;\">" + str + "<td style=\"text-align: center;font-size: 16px;height: 30px;width: 200px;border-bottom: 1px solid #999;border-right: 1px solid #999;\">" + row[i] + "</td>" + "</tr>"
		}
	}
	return info
}

//发送当前邮件
func Send(touser string, body string) {
	m := gomail.NewMessage()
	m.SetHeader("From", mainws.Li.UserInfo)
	m.SetHeader("To", touser)
	//m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", Emailconfig.Mailinfo.Title)
	m.SetBody("text/html", body)
	//m.Attach("/home/Alex/lolcat.jpg")

	iport, err := strconv.Atoi(Emailconfig.Userinfo.Port)
	checkErr(err)
	d := gomail.NewDialer(Emailconfig.Userinfo.Host, iport, Emailconfig.Userinfo.User, Emailconfig.Userinfo.Pass)

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("邮件发送成功 ！")
	}

}

//发送所有
func SendAll(xlsx string) (err error) {
	m := gomail.NewMessage()
	m.SetHeader("From", Li.UserInfo)
	m.SetHeader("To", touser)
	//m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", Emailconfig.Mailinfo.Title)
	m.SetBody("text/html", body)
	//m.Attach("/home/Alex/lolcat.jpg")

	iport, err := strconv.Atoi(Emailconfig.Userinfo.Port)
	checkErr(err)
	d := gomail.NewDialer(Emailconfig.Userinfo.Host, iport, Emailconfig.Userinfo.User, Emailconfig.Userinfo.Pass)

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("邮件发送成功 ！")
	}

}

func readXls() (xlsfile [][]string) {

	//改为Gui模式
	xlsFile, err := excelize.OpenFile()
	if err != nil {
		WarnInfo(err.Error())
	}
	sheetname := xlsFile.GetSheetName(sheetid)
	rows, err := xlsFile.GetRows(sheetname)
	if err != nil {
		WarnInfo(err.Error())
	}
	Mcell, _ = xlsFile.GetMergeCells(sheetname)
	Entry = rows[0:Head]
	Enum = len(rows) - Head
	return rows

}

func WarnInfo(str string) {
	walk.MsgBox(
		warn,
		"Error",
		str,
		walk.MsgBoxOK|walk.MsgBoxIconError)
}
