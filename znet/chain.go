package znet

import (
	"github.com/Xuzan9396/zinx-ws/ziface"
	"log"
	"reflect"
)

// base 责任链
type BaseChain struct {
	next ziface.Chain
}

func (h *BaseChain) HandleChainRequest(request ziface.ChainRequest) ziface.ChainResponse {
	if h.next != nil {
		return h.next.HandleChainRequest(request)
	}
	return request
}

func (h *BaseChain) SetNext(handler ziface.Chain) {
	h.next = handler
}

// -------------------------  全局责任链，负责设置和启动  -----------------------
type ChainAllInfo struct {
	list []ziface.Chain
}

func NewChainAllInfo() *ChainAllInfo {

	c := &ChainAllInfo{}
	// 处理消息则链路,基础责任链，默认第一条
	c.SetChain(&ChainPackMsg{})
	return c

}

func (ca *ChainAllInfo) SetChain(chain ziface.Chain) {
	if chain != nil {
		log.Println("添加责任链依次:", reflect.TypeOf(chain))
		ca.list = append(ca.list, chain)
		if len(ca.list) >= 2 {
			// 设置最新的责任链关系
			ca.list[len(ca.list)-2].SetNext(chain)
		}

	}

}

func (ca *ChainAllInfo) StartChain(request ziface.ChainRequest) ziface.ChainResponse {
	if ca.list != nil && len(ca.list) >= 0 {
		// 执行责任链触发地条
		return ca.list[0].HandleChainRequest(request)
	}
	return request
}
