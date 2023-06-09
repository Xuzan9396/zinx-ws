package ziface

// 责任链
type ChainRequest interface {
}

type ChainResErr interface {
	error
	Code() uint32
}

type Chain interface {
	HandleChainRequest(request ChainRequest) ChainResErr
	SetNext(chain Chain)
}

type ChainAll interface {
	SetChain(chain Chain)
	StartChain(request ChainRequest) ChainResErr
}
