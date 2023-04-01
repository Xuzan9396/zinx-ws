package main

import (
	wsconfig "github.com/Xuzan9396/ws/config"
	"github.com/Xuzan9396/ws/ziface"
	"github.com/Xuzan9396/ws/znet"
	"log"
	"time"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

type LoginInfo struct {
	znet.BaseRouter
}

// 模拟登录逻辑
func (l *LoginInfo) PreHandle(request ziface.IRequest) {
	auth := false
	<-time.After(5 * time.Second) // 模拟业务
	if auth == false {
		// 认证失败
		request.GetConnetion().Stop()
	}
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

// 创建链接后初始化函数
func SetOnConnetStart(conn ziface.IConnection) {
	conn.SetProperty("name", "xuzan")
	res, bools := conn.GetProperty("name")
	if bools {
		log.Println("name", res.(string))
	}
	conn.RemoveProperty("name")
}

func GetConnectNum(s ziface.IServer) {
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
}

func main() {
	// 创建一个server 句柄
	wsconfig.SetWSConfig("127.0.0.1", 8999, wsconfig.WithName("gitxuzan ----- websocket"), wsconfig.WithWorkerSize(10))
	s := znet.NewServer()
	// 启动sever
	s.SetOnConnStart(SetOnConnetStart)
	// 测试业务
	s.AddRouter(1, &HelloInfo{})
	// 其他业务
	s.AddRouter(2, &PingInfo{})
	// 登录
	s.AddRouter(1001, &LoginInfo{})

	// 监控长连接数量
	GetConnectNum(s)
	s.Server()
}
