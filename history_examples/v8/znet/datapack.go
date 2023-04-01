package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/Xuzan9396/ws/v8/utils"
	"github.com/Xuzan9396/ws/v8/ziface"
	"log"
)

type DataPack struct {
}

func NewDataPack() *DataPack {
	return &DataPack{}
}

// 封包
func (c *DataPack) Pack(msg ziface.Imessage) ([]byte, error) {
	// 创建字节流缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	// datalen 写入 dataBuff中
	if err := binary.Write(dataBuff, binary.BigEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	// msg id
	if err := binary.Write(dataBuff, binary.BigEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	// data写入数据中
	if err := binary.Write(dataBuff, binary.BigEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

// 拆包
func (c *DataPack) Unpack(b []byte) (ziface.Imessage, error) {
	dataBuff := bytes.NewReader(b)

	msg := &Message{}

	if err := binary.Read(dataBuff, binary.BigEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	log.Println("--------------------")
	//log.Println(msg.DataLen, utils.GlobalObject.MaxPackageSize)
	if msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("接收的包超出范围了!")
	}

	if err := binary.Read(dataBuff, binary.BigEndian, &msg.Id); err != nil {
		return nil, err
	}
	msg.Data = make([]byte, msg.DataLen)
	if err := binary.Read(dataBuff, binary.BigEndian, &msg.Data); err != nil {
		return nil, err
	}
	//msg.Data = b[8 : 8+msg.DataLen]

	return msg, nil

}
