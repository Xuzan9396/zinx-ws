package main

import (
	"github.com/Xuzan9396/ws/v9/ziface"
	"github.com/Xuzan9396/ws/v9/znet"
	"log"
	"time"
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
	s := znet.NewServer("[zinx-ws v9]")
	// 启动sever
	s.SetOnConnStart(func(conn ziface.IConnection) {

		conn.SetProperty("name", "xuzan")
		res, bools := conn.GetProperty("name")
		if bools {
			log.Println("name", res.(string))
		}
		conn.RemoveProperty("name")

	})
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				connNumTotal := s.GetConnMgr().Len()
				log.Println("连接数量:", connNumTotal)
			}
		}
	}()

	s.AddRouter(1, &HelloInfo{})
	s.AddRouter(2, &PingInfo{})
	s.Server()
}
