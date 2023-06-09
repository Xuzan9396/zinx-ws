package ziface

/*
封包 拆包
*/

type IDataPack interface {

	// 封包
	Pack(msg Imessage) ([]byte, error)

	// 拆包
	Unpack([]byte) (Imessage, error)
}
