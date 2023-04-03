package main

import (
	"github.com/Xuzan9396/zinx-ws/v5/ziface"
	"github.com/Xuzan9396/zinx-ws/v5/znet"
	"log"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

type PingInfo struct {
	znet.BaseRouter
}

func (p *PingInfo) PreHandle(request ziface.IRequest) {
	log.Printf("pre:%s,conntId:%d,msgId:%d", request.GetData(), request.GetConnetion().GetConnID(), request.GetMsgId())
}

func (p *PingInfo) Handle(request ziface.IRequest) {
	log.Printf("Handle:%s,conntId:%d,msgId:%d", request.GetData(), request.GetConnetion().GetConnID(), request.GetMsgId())

}

func (p *PingInfo) PostHandle(request ziface.IRequest) {
	log.Printf("post:%s,conntId:%d,,msgId:%d", request.GetData(), request.GetConnetion().GetConnID(), request.GetMsgId())

}

func main() {
	// 创建一个server 句柄
	s := znet.NewServer("[zinx-ws v5]")
	// 启动sever
	s.AddRouter(&PingInfo{})
	s.Server()
}
