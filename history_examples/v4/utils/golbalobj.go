package utils

import (
	"encoding/json"
	"github.com/Xuzan9396/ws/v4/ziface"
	"io/ioutil"
)

type GlobalObj struct {
	TcpServer ziface.IServer // tcp server对象
	Host      string
	TcpPort   int
	Name      string //

	Version        string
	MaxConn        int    // 最大连接数
	MaxPackageSize uint32 // 最大数据包大小
}

var GlobalObject *GlobalObj

func init() {
	GlobalObject = &GlobalObj{
		Name:           "jiajia",
		Version:        "0.4",
		TcpPort:        8999,
		Host:           "0.0.0.0",
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}
	GlobalObject.Reload()
	//for i, i2 := range GlobalObject {
	//
	//}
}

func (g *GlobalObj) Reload() {
	byteS, err := ioutil.ReadFile("zinx.json")
	if err != nil {
		//log.Fatal(err)
		//log.Println(err)
		return
	}

	json.Unmarshal(byteS, GlobalObject)
}
