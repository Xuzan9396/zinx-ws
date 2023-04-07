package ziface

type IRequest interface {
	// 当前链接

	GetConnetion() IConnection
	// 得到请求消息数据
	GetData() []byte

	GetMsgID() uint32
}
