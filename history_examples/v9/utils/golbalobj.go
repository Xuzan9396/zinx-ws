package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type GlobalObj struct {
	Host           string
	TcpPort        int
	Name           string
	Version        string
	MaxConn        int    // 最大连接数
	MaxPackageSize uint32 // 最大数据包大小
	WorkerPoolSize uint32 // 工作池大小
}

var GlobalObject *GlobalObj

func init() {
	GlobalObject = &GlobalObj{
		Name:           "zinx-ws",
		Version:        "v9",
		TcpPort:        8999,
		Host:           "127.0.0.1",
		MaxConn:        1000,
		MaxPackageSize: 4096,
		WorkerPoolSize: 1,
	}
	GlobalObject.Reload("../zinx.json")

}

func (g *GlobalObj) Reload(fileName string) {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		// 文件不存在
		log.Println(err)
		return
	}

	byteS, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Println(err)
		return
	}

	if err = json.Unmarshal(byteS, GlobalObject); err != nil {
		log.Println(err)
		return
	}
}
