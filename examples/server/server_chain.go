package main

import (
	wsconfig "github.com/Xuzan9396/zinx-ws/config"
	"github.com/Xuzan9396/zinx-ws/ziface"
	"github.com/Xuzan9396/zinx-ws/znet"
	"log"
	"time"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

type MyChain struct {
	znet.BaseChain
}

func (h *MyChain) HandleChainRequest(request ziface.ChainRequest) ziface.ChainResErr {
	log.Println("我是责任链第二条!")
	time.Sleep(5 * time.Second)
	return znet.NewError("错误sss", 10002)
	//return h.BaseChain.HandleChainRequest(request)
}

func main() {
	// 创建一个server 句柄
	wsconfig.SetWSConfig("127.0.0.1", 8999, wsconfig.WithName("gitxuzan ----- websocket"), wsconfig.WithWorkerSize(10))
	s := znet.NewServer()

	s.AddChain(new(MyChain))

	s.Server()
}
