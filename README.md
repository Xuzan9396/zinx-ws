# zinx-websocket版本  
#### 目标?
##### 维护golang版websocket版本, 打算跟zinx tcp 版本同步，然后无偿开源

更新
【v0.0.2】新增责任链，新增自定义责任链,消息拦截解耦,examples有案例,，添加znet/chain.go , ziface/ichain.go


#### 为什么做这个项目？
`
做这个项目初衷，主要因为自己公司做直播平台的，之前公司写了一套，websocket封装的框架，主要做房间服务器,和h5小游戏服务器，但是由于感觉随着业务增大，后面感觉某些设计有缺陷，看了冰哥的设计模式， 打算跟着冰哥设计模式重写一个websocket
`

#### 打算项目使用？

```
后续会在自己的项目中使用，打算在直播间的小游戏，准备上线使用
```


#### 具体怎么使用

##### 参数说明
```json

  "Name": "zin-ws -------gitxuzan",
  "Host": "127.0.0.1",
  "端口": "端口",
  "TcpPort": 8999,
  "最大连接数": "最大连接数",
  "MaxConn": 1000,
  "最大的包大小": "最大包大小",
  "MaxPackageSize": 4096,
  "worker池子": "worker池子10个并发处理读的数据",
  "WorkerPoolSize": 10

```

##### 数据发送格式简单说明(后续修改成格式定义)
| MsgId  | len    | body      |
|--------|--------|-----------|
| 协议号ID  | body长度 | 二进制body长度 | 
| uint32 | uint32 | []byte    | 


##### 服务端配置设置
```golang
wsconfig.SetWSConfig("127.0.0.1", 8999, wsconfig.WithName("gitxuzan ----- websocket"))
还有其他设置例如:
wsconfig.WithWorkerSize(10) // 设置10个worker处理业务逻辑
wsconfig.WithMaxPackSize(4096)  // 每个发送的包大小 4k
wsconfig.WithMaxConn(1000)	// 同时在线1000个连接     
wsconfig.WithVersion()	        // 自定义本地版本   
```

##### 定义业务逻辑协议
```golang 
type LoginInfo struct {
	znet.BaseRouter
}

例如上面写的 LoginInfo 继承znet.BaseRouter
重写三个方法依次执行:
PreHandle
Handle
PostHandle
```
##### 设置router 映射到具体的方法上
```golang
同时要设置router 
	// 登录
s.AddRouter(1001, &LoginInfo{})
1001 代表协议号，相当于协议投里面的msgId,映射到具体某个业务，发送端需要发送对应的协议号
```

##### request 的一些功能，例如下面的案例，模拟登入验证等等
```
func (l *LoginInfo) PreHandle(request ziface.IRequest) {
request 中 目前有发送，断开，获取当前属性，获取当前连接
}
```



##### 完整的服务端使用代码
```golang
package main

import (
	wsconfig "github.com/Xuzan9396/zinx-ws/config"
	"github.com/Xuzan9396/zinx-ws/ziface"
	"github.com/Xuzan9396/zinx-ws/znet"
	"log"
	"time"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

type LoginInfo struct {
	znet.BaseRouter
}

// 模拟登录逻辑
func (l *LoginInfo) PreHandle(request ziface.IRequest) {
	auth := false
	<-time.After(5 * time.Second) // 模拟业务
	if auth == false {
		// 模拟登录认证失败，然后断开连接
		request.GetConnetion().Stop()
	}
}

type PingInfo struct {
	znet.BaseRouter
}
type HelloInfo struct {
	znet.BaseRouter
}

func (p *PingInfo) PreHandle(request ziface.IRequest) {
	log.Printf("pre:%s,conntId:%d,msgId:%d", request.GetData(), request.GetConnetion().GetConnID(), request.GetMsgID())
}

func (p *PingInfo) Handle(request ziface.IRequest) {
	log.Printf("Handle:%s,conntId:%d,msgId:%d", request.GetData(), request.GetConnetion().GetConnID(), request.GetMsgID())

}

func (p *PingInfo) PostHandle(request ziface.IRequest) {
	log.Printf("post:%s,conntId:%d,,msgId:%d", request.GetData(), request.GetConnetion().GetConnID(), request.GetMsgID())
	request.GetConnetion().SendMsg(request.GetMsgID(), []byte("回复ping!"))
}
func (p *HelloInfo) PreHandle(request ziface.IRequest) {
	log.Printf("pre:%s,conntId:%d,msgId:%d", request.GetData(), request.GetConnetion().GetConnID(), request.GetMsgID())
}

func (p *HelloInfo) Handle(request ziface.IRequest) {
	log.Printf("Handle:%s,conntId:%d,msgId:%d", request.GetData(), request.GetConnetion().GetConnID(), request.GetMsgID())

}

func (p *HelloInfo) PostHandle(request ziface.IRequest) {
	log.Printf("post:%s,conntId:%d,,msgId:%d", request.GetData(), request.GetConnetion().GetConnID(), request.GetMsgID())
	request.GetConnetion().SendMsg(request.GetMsgID(), []byte("回复hello!"))

}

// 创建链接后初始化函数
func SetOnConnetStart(conn ziface.IConnection) {
	conn.SetProperty("name", "xuzan")
	res, bools := conn.GetProperty("name")
	if bools {
		log.Println("name", res.(string))
	}
	conn.RemoveProperty("name")
}

func GetConnectNum(s ziface.IServer) {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				connNumTotal := s.GetConnMgr().Len()
				log.Println("连接数量:", connNumTotal)
			}
		}
	}()
}

func main() {
	//设置配置
	wsconfig.SetWSConfig("127.0.0.1", 8999, wsconfig.WithName("gitxuzan ----- websocket"))
	// 创建一个server 句柄
	s := znet.NewServer()
	// 启动sever
	s.SetOnConnStart(SetOnConnetStart)
	// 测试业务
	s.AddRouter(1, &HelloInfo{})
	// 其他业务
	s.AddRouter(2, &PingInfo{})
	// 登录
	s.AddRouter(1001, &LoginInfo{})

	// 监控长连接数量
	GetConnectNum(s)
	s.Server()
}


```


##### 完整的客户端案例代码
```golang
package main

import (
	"flag"
	"github.com/Xuzan9396/zinx-ws/znet"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"
)

var addr = flag.String("addr", "127.0.0.1:8999", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/"}
	log.Printf("connecting to %s", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), http.Header{"User-Agent": {""}})
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	log.Println("ws 连接成功")
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	p := znet.NewDataPack()
	by := []byte{'h', 'e', 'l', 'l', 'o'}
	resBytes, err := p.Pack(&znet.Message{
		Id:      1,
		DataLen: uint32(len(by)),
		Data:    by,
	})

	byPing := []byte("ping")
	resPingBytes, _ := p.Pack(&znet.Message{
		Id:      2,
		DataLen: uint32(len(byPing)),
		Data:    byPing,
	})
	timer := time.NewTimer(30 * time.Second)

	go read(c)
	for {
		select {
		case <-timer.C:
			// 模拟认证登录
			sendMsg := []byte("login")
			sendMsgPack, _ := p.Pack(&znet.Message{
				Id:      1001,
				DataLen: uint32(len(sendMsg)),
				Data:    sendMsg,
			})
			err := c.WriteMessage(websocket.BinaryMessage, sendMsgPack)
			if err != nil {
				log.Println("write:", err)
				timer.Stop()
				return
			}
			log.Println("login写入成功:", string(sendMsg))
			timer.Stop()
		case <-ticker.C:
			sendMsg := resBytes
			err := c.WriteMessage(websocket.BinaryMessage, sendMsg)
			if err != nil {
				log.Println("write:", err)
				return
			}
			log.Println("写入成功:", string(by))

			err = c.WriteMessage(websocket.BinaryMessage, resPingBytes)
			if err != nil {
				log.Println("write:", err)
				return
			}

			log.Println("写入成功:", string(resPingBytes))
		case <-interrupt:
			log.Println("interrupt")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}

		}
	}

}

func read(c *websocket.Conn) {
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		p := znet.NewDataPack()
		img, err := p.Unpack(message)
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("msgId:%d,recv: %s", img.GetMsgId(), img.GetData())
	}
}

```

###### 参考链接
https://github.com/aceld/zinx 来自冰哥

[![Visitors](https://visitor-badge.glitch.me/badge?page_id=Xuzan9396.zinx-ws)](https://github.com/Xuzan9396/zinx-ws)
