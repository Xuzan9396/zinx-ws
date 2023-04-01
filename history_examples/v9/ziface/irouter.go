package ziface

// 路由抽象接口
// 路由数据都是iReqest

type IRouter interface {
	PreHandle(request IRequest)

	// 处理中
	Handle(request IRequest)

	// 处理后
	PostHandle(request IRequest)
}