package ziface

type IServer interface {
	Start()
	Stop()
	Server()
	AddRouter(msgId uint32, router IRouter)
}
