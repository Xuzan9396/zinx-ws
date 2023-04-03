package znet

import (
	"errors"
	wsconfig "github.com/Xuzan9396/zinx-ws/config"
	"github.com/Xuzan9396/zinx-ws/ziface"

	"github.com/gorilla/websocket"
	"log"
	"net"
	"runtime/debug"
	"sync"
	"time"
)

type Connection struct {
	// 当前连接socket tcp 套接字
	Conn *websocket.Conn

	mutex sync.RWMutex

	// 链接的id
	ConnID uint32

	// 当前的链接状态
	isClosed bool

	// 对应msgHandler
	msgHandler ziface.ImsgHandle

	// 用户无缓冲管道 用户读写 goroutine通信
	msgChan chan []byte

	// 告知当前链接已经退出的 channel
	ExitChan chan struct{}

	// 链接管理模块
	ConnMgr ziface.IConnManager

	//当前连接创建时Hook函数
	onConnStart func(conn ziface.IConnection)

	//当前连接断开时的Hook函数
	onConnStop func(conn ziface.IConnection)

	// 属性
	property map[string]interface{}

	propertyLock sync.RWMutex
}

const (
	writeWait  = 10 * time.Second // 写入超时时间
	pongWait   = 60 * time.Second //读取设置超时 reader端 pongWait 要比 pingPeriod大
	pingPeriod = 9 * time.Second  // writer 端定时发送ping
)

// 初始化
func NewConnection(server ziface.IServer, conn *websocket.Conn, connID uint32) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		msgChan:  make(chan []byte, 256),
		ExitChan: make(chan struct{}),
	}
	// 连接管理器
	c.ConnMgr = server.GetConnMgr()
	// 消息处理器
	c.msgHandler = server.GetMsgHandler()

	c.onConnStart = server.GetOnConnStart()
	c.onConnStop = server.GetOnConnStop()

	return c
}

func (c *Connection) Start() {
	// 添加链接
	c.ConnMgr.Add(c)
	log.Println("start read and writer connect_id:", c.ConnID)
	c.callOnConnStart()
	// 启动读服务
	go c.StartReader()
	// 启动写服务
	go c.StartWriter()
}

func (c *Connection) StartReader() {

	defer log.Println("connid = ", c.ConnID, " reader is exit ")
	defer c.Stop()
	c.Conn.SetReadLimit(int64(wsconfig.GlobalObject.MaxPackageSize))
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(
		func(string) error {
			c.Conn.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})

	for {
		message, err := c.read()
		if err != nil {
			break
		}
		dp := NewDataPack()
		imsg, err := dp.Unpack(message)
		if err != nil {
			log.Printf("error: %v\n", err)
			break

		}
		req := &Request{
			conn: c,
			msg:  imsg,
		}
		//c.msgHandler.DoMsgHandler(req) // 没workerpool
		c.msgHandler.SendMsgToTaskQueue(req) // 设置workerpool分发

	}
}

func (c *Connection) read() ([]byte, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if c.isClosed {
		return nil, errors.New("conn close,reader")
	}
	_, message, err := c.Conn.ReadMessage()
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			log.Printf("error: %v", err)
		} else {
			log.Println("其他", err)
		}
		return nil, err
	}
	return message, err
}

func (c *Connection) StartWriter() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		if err := recover(); err != nil {
			log.Println(string(debug.Stack()))
		}
	}()

	for {
		select {

		case data, ok := <-c.msgChan:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			c.write(websocket.BinaryMessage, data)

		case <-c.ExitChan:
			log.Println("conneid = ", c.ConnID, ",write 退出")
			return
		case <-ticker.C:
			c.write(websocket.PingMessage, nil)

		}
	}
}
func (c *Connection) write(messageType int, data []byte) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if c.isClosed == true {
		log.Println(errors.New("connection close,连接已经关闭了 "))
	}
	c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
	c.Conn.WriteMessage(messageType, data)
}

func (c *Connection) Stop() {
	c.callOnConnStop()
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.isClosed {
		return
	}
	c.isClosed = true
	c.Conn.Close()
	c.ConnMgr.Remove(c)
	close(c.ExitChan)
	close(c.msgChan)

}

// 获取wss socket
func (c *Connection) GetTCPConnection() *websocket.Conn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 发送给客户端
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if c.isClosed == true {
		return errors.New("连接已关闭!")
	}
	dp := NewDataPack()

	msgByte, err := dp.Pack(&Message{
		Id:      msgId,
		DataLen: uint32(len(data)),
		Data:    data,
	})

	if err != nil {
		return err
	}
	c.msgChan <- msgByte
	return nil
}

// callOnConnStart 调用连接OnConnStart Hook函数
func (c *Connection) callOnConnStart() {
	if c.onConnStart != nil {
		log.Println("CallOnConnStart....")
		c.onConnStart(c)
	}
}

// callOnConnStop 调用连接OnConnStop Hook函数
func (c *Connection) callOnConnStop() {
	if c.onConnStop != nil {
		log.Println("CallOnConnStop....")
		c.onConnStop(c)
	}
}

// SetProperty 设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	if c.property == nil {
		c.property = make(map[string]interface{})
	}
	c.property[key] = value
}

// GetProperty 获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, bool) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, ok
	}
	return nil, false
}

// RemoveProperty 移除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property, key)
}
