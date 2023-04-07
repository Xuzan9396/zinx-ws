package znet

import (
	"github.com/Xuzan9396/zinx-ws/ziface"
)

// 处理消息的责链路
type ChainPackMsg struct {
	BaseChain
}

func (h *ChainPackMsg) HandleChainRequest(request ziface.ChainRequest) (res ziface.ChainResponse) {

	ireq, ok := request.(ziface.IRequest)
	if !ok {
		return nil
	}
	ireq.GetConnetion().GetmsgHandler().SendMsgToTaskQueue(ireq) // 设置workerpool分发

	return h.BaseChain.HandleChainRequest(request)
}
