package znet

import (
	"github.com/Xuzan9396/zinx-ws/v7/ziface"
	"log"
	"sync"
)

type MsgHandle struct {
	Apis map[uint32]ziface.IRouter
	sync.Mutex
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),
	}
}

func (c *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	handler, ok := c.Apis[request.GetMsgID()]
	if !ok {
		log.Println("没有找到该方法")
		return
	}
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (c *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) {
	c.Lock()
	defer c.Unlock()
	c.Apis[msgId] = router
}
