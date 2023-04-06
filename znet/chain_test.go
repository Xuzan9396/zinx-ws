package znet

import (
	"github.com/Xuzan9396/zinx-ws/ziface"
	"log"
	"testing"
)

type HandlerChainA struct {
	BaseChain
}

func (h *HandlerChainA) HandleChainRequest(request ziface.ChainRequest) ziface.ChainResponse {
	log.Println("执行了A!!")
	// 执行b
	return h.BaseChain.HandleChainRequest(request)
}

type HandlerChainB struct {
	BaseChain
}

func (h *HandlerChainB) HandleChainRequest(request ziface.ChainRequest) ziface.ChainResponse {
	log.Println("执行了B!!")
	return h.BaseChain.HandleChainRequest(request)
}

type HandlerChainC struct {
	BaseChain
}

func (h *HandlerChainC) HandleChainRequest(request ziface.ChainRequest) ziface.ChainResponse {
	log.Println("执行了C!!")
	return h.BaseChain.HandleChainRequest(request)
}

func Test_chain(t *testing.T) {
	handlerA := &HandlerChainA{}
	handlerB := &HandlerChainB{}
	handlerC := &HandlerChainC{}
	handlerA.SetNext(handlerB)
	handlerB.SetNext(handlerC)

	// 触发 依次 打印 A -> B -> C
	handlerA.HandleChainRequest("")
}
