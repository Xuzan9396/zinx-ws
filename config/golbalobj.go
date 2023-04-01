package config

type GlobalObj struct {
	Host           string
	TcpPort        int
	Name           string
	Version        string
	MaxConn        int    // 最大连接数
	MaxPackageSize uint32 // 最大数据包大小
	WorkerPoolSize uint32 // 工作池大小
}

var GlobalObject = &GlobalObj{
	Name:           "zinx-ws",
	Version:        "v0.0.1",
	TcpPort:        8999,
	Host:           "127.0.0.1",
	MaxConn:        1000,
	MaxPackageSize: 4096,
	WorkerPoolSize: 2,
}

type ServerOption func(*GlobalObj)

func WithName(name string) ServerOption {
	return func(obj *GlobalObj) {
		obj.Name = name
	}
}

func WithVersion(version string) ServerOption {
	return func(obj *GlobalObj) {
		obj.Version = version
	}
}

func WithMaxConn(maxconn int) ServerOption {
	return func(obj *GlobalObj) {
		obj.MaxConn = maxconn
	}
}

func WithMaxPackSize(maxPackageSize uint32) ServerOption {
	return func(obj *GlobalObj) {
		obj.MaxPackageSize = maxPackageSize
	}
}

func WithWorkerSize(workerNum uint32) ServerOption {
	return func(obj *GlobalObj) {
		obj.WorkerPoolSize = workerNum
	}
}

func SetWSConfig(host string, port int, options ...ServerOption) {
	GlobalObject.Host = host
	GlobalObject.TcpPort = port

	for _, option := range options {
		option(GlobalObject)
	}

}
