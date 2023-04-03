package znet

import "github.com/Xuzan9396/zinx-ws/v8/ziface"

type Request struct {
	// 已经和客户端建立好链接
	conn ziface.IConnection
	msg  ziface.Imessage
}

func (r *Request) GetConnetion() ziface.IConnection {
	return r.conn
}

// 获取数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgId()
}
