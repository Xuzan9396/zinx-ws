package ziface

type Imessage interface {

	// 获取消息的ID
	GetMsgId() uint32
	// 获取消息的长度
	GetMsgLen() uint32
	// 获取消息的内容
	GetData() []byte
	// 设置消息的ID
	SetMsgId(uint32)

	SetData([]byte)

	SetDataLen(uint32)
}
