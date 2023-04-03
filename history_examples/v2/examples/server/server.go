package main

import (
	"github.com/Xuzan9396/zinx-ws/v2/znet"
	"log"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

}
func main() {
	// 创建一个server 句柄
	s := znet.NewServer("[zinx-ws v0.2]")
	// 启动sever
	s.Server()
}
