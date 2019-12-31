package mainws

import (
	"AutoSalaryGui/loginws"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const ConfigName = "autosalary.config"

func SaveLogin(li *loginws.LoginInfo) {

	//	var filePtr *os.File
	passbyte := []byte(li.PassInfo)
	//加密
	encoded := base64.StdEncoding.EncodeToString(passbyte)
	//解密
	//decoded, _ := base64.StdEncoding.DecodeString(encoded)
	li.PassInfo = encoded
	data, _ := json.Marshal(li)

	/*	path, _ := os.Getwd()
		fmt.Println(path)
		filepath := path + "/" + ConfigName
		if !FileExist(filepath) {
			filePtr, _ = os.Create(ConfigName)
		} else {
			filePtr, _ = os.OpenFile(ConfigName, os.O_TRUNC, 0660)
		}*/
	err := ioutil.WriteFile(ConfigName, data, 0660)
	if err != nil {
		fmt.Print("保存用户信息失败！")
	}
}

func ReadConf(li *loginws.LoginInfo) {
	path, _ := os.Getwd()
	filepath := path + "/" + ConfigName
	if FileExist(filepath) {
		filePtr, err := os.Open(ConfigName)
		if err != nil {
			fmt.Println("读取用户信息失败！")
		}
		defer filePtr.Close()
		readInfo := json.NewDecoder(filePtr)
		err = readInfo.Decode(li)
		if err != nil {
			fmt.Println("存储用户信息格式错误！")
		}
		decoded, _ := base64.StdEncoding.DecodeString(li.PassInfo)
		li.PassInfo = string(decoded)
	}
}

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func ChooseFile() {

}

func Reset() {

}
