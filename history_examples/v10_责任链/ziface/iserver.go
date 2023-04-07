package ziface

type IServer interface {
	Start()
	Stop()
	Server()
	AddRouter(msgId uint32, router IRouter)
	GetConnMgr() IConnManager
	GetMsgHandler() ImsgHandle
	SetOnConnStart(hookFunc func(conn IConnection))
	SetOnConnStop(hookFunc func(conn IConnection))
	GetOnConnStart() (hookFunc func(conn IConnection))
	GetOnConnStop() (hookFunc func(conn IConnection))

	AddChain(chain Chain) // 设置责任链
	GetChain() ChainAll   // 设置责任链

}
