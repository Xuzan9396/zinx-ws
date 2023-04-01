package ziface

import (
	"github.com/gorilla/websocket"
	"net"
)

type IConnection interface {

	// 启动链接
	Start()
	// 停止链接
	Stop()
	// 获取当前链接绑定 socket conn
	GetTCPConnection() *websocket.Conn

	// 获取当前链接模块链接ID
	GetConnID() uint32

	// 获取原创客户端tcp状态
	RemoteAddr() net.Addr

	// 发送数据，给远程客户端
	Send(data []byte) error
}

type HandleFunc func(*websocket.Conn, []byte, int) error
