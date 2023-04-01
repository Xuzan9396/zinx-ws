package znet

import (
	"testing"
)

func TestDataPack_UnPack(t *testing.T) {
	p := NewDataPack()

	by := []byte{'h', 'e', 'l', 'l', 'o'}
	resBytes, err := p.Pack(&Message{
		Id:      1,
		DataLen: uint32(len(by)),
		Data:    by,
	})
	if err != nil {
		t.Log(err)
		return
	}
	t.Log("封包成功!")

	iMessage, err := p.Unpack(resBytes)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log("解包成功", iMessage.GetMsgLen(), iMessage.GetMsgId(), string(iMessage.GetData()))
}
