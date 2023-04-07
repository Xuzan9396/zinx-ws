package ziface

/*
封包 拆包
*/

type IDataPack interface {
	// 长度
	GetHeadLen() uint32

	// 封包
	Pack(msg Imessage) ([]byte, error)

	// 拆包
	Unpack([]byte) (Imessage, error)
}
