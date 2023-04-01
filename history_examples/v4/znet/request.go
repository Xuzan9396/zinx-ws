package znet

import "github.com/Xuzan9396/ws/v4/ziface"

type Request struct {
	// 已经和客户端建立好链接
	conn ziface.IConnection
	data []byte
}

func (r *Request) GetConnetion() ziface.IConnection {
	return r.conn
}

// 获取数据
func (r *Request) GetData() []byte {
	return r.data
}
