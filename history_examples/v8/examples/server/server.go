package main

import (
	"github.com/Xuzan9396/zinx-ws/v8/ziface"
	"github.com/Xuzan9396/zinx-ws/v8/znet"
	"log"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

type PingInfo struct {
	znet.BaseRouter
}
type HelloInfo struct {
	znet.BaseRouter
}

func (p *PingInfo) PreHandle(request ziface.IRequest) {
	log.Printf("pre:%s,conntId:%d,msgId:%d", request.GetData(), request.GetConnetion().GetConnID(), request.GetMsgID())
}

func (p *PingInfo) Handle(request ziface.IRequest) {
	log.Printf("Handle:%s,conntId:%d,msgId:%d", request.GetData(), request.GetConnetion().GetConnID(), request.GetMsgID())

}

func (p *PingInfo) PostHandle(request ziface.IRequest) {
	log.Printf("post:%s,conntId:%d,,msgId:%d", request.GetData(), request.GetConnetion().GetConnID(), request.GetMsgID())
	request.GetConnetion().SendMsg(request.GetMsgID(), []byte("回复ping!"))
}
func (p *HelloInfo) PreHandle(request ziface.IRequest) {
	log.Printf("pre:%s,conntId:%d,msgId:%d", request.GetData(), request.GetConnetion().GetConnID(), request.GetMsgID())
}

func (p *HelloInfo) Handle(request ziface.IRequest) {
	log.Printf("Handle:%s,conntId:%d,msgId:%d", request.GetData(), request.GetConnetion().GetConnID(), request.GetMsgID())

}

func (p *HelloInfo) PostHandle(request ziface.IRequest) {
	log.Printf("post:%s,conntId:%d,,msgId:%d", request.GetData(), request.GetConnetion().GetConnID(), request.GetMsgID())
	request.GetConnetion().SendMsg(request.GetMsgID(), []byte("回复hello!"))

}

func main() {
	// 创建一个server 句柄
	s := znet.NewServer("[zinx-ws v8]")
	// 启动sever
	s.AddRouter(1, &HelloInfo{})
	s.AddRouter(2, &PingInfo{})
	s.Server()
}
