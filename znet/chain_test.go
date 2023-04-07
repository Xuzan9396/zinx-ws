package znet

import (
	"github.com/Xuzan9396/zinx-ws/ziface"
	"log"
	"testing"
)

type HandlerChainA struct {
	BaseChain
}

func (h *HandlerChainA) HandleChainRequest(request ziface.ChainRequest) ziface.ChainResErr {
	log.Println("执行了A!!")
	// 执行b
	return h.BaseChain.HandleChainRequest(request)
}

type HandlerChainB struct {
	BaseChain
}

func (h *HandlerChainB) HandleChainRequest(request ziface.ChainRequest) ziface.ChainResErr {
	log.Println("执行了B!!")
	return h.BaseChain.HandleChainRequest(request)
}

type HandlerChainC struct {
	BaseChain
}

func (h *HandlerChainC) HandleChainRequest(request ziface.ChainRequest) ziface.ChainResErr {
	log.Println("执行了C!!")
	return h.BaseChain.HandleChainRequest(request)
}

type HandlerChainD struct {
	BaseChain
}

func (h *HandlerChainD) HandleChainRequest(request ziface.ChainRequest) ziface.ChainResErr {
	log.Println("执行了D,但是报错了!!")

	return NewError("我是D,我报错了", 1000)
	//return h.BaseChain.HandleChainRequest(request)
}

func Test_chain(t *testing.T) {
	handlerA := &HandlerChainA{}
	handlerB := &HandlerChainB{}
	handlerC := &HandlerChainC{}
	handlerD := &HandlerChainD{}
	handlerA.SetNext(handlerB)
	handlerB.SetNext(handlerC)
	handlerC.SetNext(handlerD)

	// 触发 依次 打印 A -> B -> C -> D
	err := handlerA.HandleChainRequest("xxxx")
	if err != nil {
		log.Println(err.Error())
	}
}

func Test_chainZhongduan(t *testing.T) {
	handlerA := &HandlerChainA{}
	handlerB := &HandlerChainB{}
	handlerC := &HandlerChainC{}
	handlerD := &HandlerChainD{}
	handlerA.SetNext(handlerB)
	handlerB.SetNext(handlerD)
	handlerD.SetNext(handlerC)

	// 触发 依次 打印 A -> B -> D 终端 -> 不会执行C
	err := handlerA.HandleChainRequest("xxxx")
	if err != nil {
		log.Println(err.Error())
	}
}
