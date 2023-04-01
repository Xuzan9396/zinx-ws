package ziface

type ImsgHandle interface {
	DoMsgHandler(request IRequest)

	AddRouter(msgId uint32,router IRouter)
}