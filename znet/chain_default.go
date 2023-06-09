package znet

import (
	"github.com/Xuzan9396/zinx-ws/ziface"
)

// 处理消息的责链路
type ChainPackMsg struct {
	BaseChain
}

// 责任链第一条
func (h *ChainPackMsg) HandleChainRequest(request ziface.ChainRequest) (res ziface.ChainResErr) {
	ireq, ok := request.(ziface.IRequest)
	if !ok {
		return NewError("ziface.IRequest 断言失败", 1000)
	}

	ireq.GetConnetion().GetmsgHandler().SendMsgToTaskQueue(ireq) // 设置workerpool分发

	return h.BaseChain.HandleChainRequest(request)
}
