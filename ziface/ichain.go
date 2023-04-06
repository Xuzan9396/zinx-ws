package ziface

// 责任链
type ChainRequest interface {
}

type ChainResponse interface {
}
type Chain interface {
	HandleChainRequest(request ChainRequest) ChainResponse
	SetNext(chain Chain)
}

type ChainAll interface {
	SetChain(chain Chain)
	StartChain(request ChainRequest) ChainResponse
}
