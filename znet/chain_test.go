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
	arr := []ziface.Chain{handlerA, handlerB, handlerC} // 按照这个顺序执行
	lens := len(arr)
	if lens > 0 {
		for i, chain := range arr {
			// 触发运行第一个链条，执行A 然后执行b
			if i+1 < lens {
				chain.SetNext(arr[i+1])
			}
		}
		// 触发运行第一个链条，执行A 然后执行b , 再执行c
		arr[0].HandleChainRequest(nil)
	}

}
