package znet

import "github.com/Xuzan9396/zinx-ws/ziface"

type Request struct {
	// 已经和客户端建立好链接
	conn ziface.IConnection
	msg  ziface.Imessage
}

func NewRequest(c ziface.IConnection, msg ziface.Imessage, raw_data []byte) *Request {
	return &Request{
		conn: c,
		msg:  msg,
	}
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
