package ziface

type ImsgHandle interface {
	// 单个发送
	DoMsgHandler(request IRequest)

	// 多个woker池轮询发送
	SendMsgToTaskQueue(request IRequest)

	AddRouter(msgId uint32, router IRouter)
}
