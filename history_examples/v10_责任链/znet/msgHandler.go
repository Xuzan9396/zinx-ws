package znet

import (
	wsconfig "github.com/Xuzan9396/zinx-ws/config"
	"github.com/Xuzan9396/zinx-ws/ziface"

	"log"
	"runtime/debug"
	"sync"
)

type MsgHandle struct {
	Apis map[uint32]ziface.IRouter
	// 负责worker取任务的消息队列
	TaskQueue []chan ziface.IRequest
	// worker 工作池的数量
	WorkerPoolSize uint32
	sync.Mutex
}

func NewMsgHandle() *MsgHandle {

	msgHandle := &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: wsconfig.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, wsconfig.GlobalObject.WorkerPoolSize),
	}
	msgHandle.StartWorkerPool()
	return msgHandle

}

// 启动一个worker工作池 启动一次
func (c *MsgHandle) StartWorkerPool() {
	for i := 0; i < int(c.WorkerPoolSize); i++ {
		// 对应worker开闭管道
		c.TaskQueue[i] = make(chan ziface.IRequest, 1024)
		// 启动当前worker,阻塞消息
		go c.StartOneWorker(i, c.TaskQueue[i])
	}
}

func (c *MsgHandle) StartOneWorker(i int, req chan ziface.IRequest) {
	for {
		select {
		case request := <-req:
			c.DoMsgHandler(request)
		}
	}
}

func (c *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	defer func() {
		if err := recover(); err != nil {
			// 防止panic
			log.Println(string(debug.Stack()))
		}
	}()
	handler, ok := c.Apis[request.GetMsgID()]
	if !ok {
		log.Println("没有找到该方法的协议号")
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

// 消息给taskQueue 平均分配
func (c *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	workerID := request.GetConnetion().GetConnID() % c.WorkerPoolSize
	c.TaskQueue[workerID] <- request
}
