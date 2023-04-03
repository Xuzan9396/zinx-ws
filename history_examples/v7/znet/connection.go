package znet

import (
	"errors"
	"github.com/Xuzan9396/zinx-ws/v7/utils"
	"github.com/Xuzan9396/zinx-ws/v7/ziface"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"runtime/debug"
	"time"
)

type Connection struct {
	// 当前连接诶socket tcp 套接字

	Conn *websocket.Conn

	// 链接的id
	ConnID uint32

	// 当前的链接状态
	isClosed bool

	// 对应msgHandler
	msgHandler ziface.ImsgHandle

	// 用户无缓冲管道 用户读写 goroutine通信
	msgChan chan []byte

	// 告知当前链接已经退出的 channel
	ExitChan chan bool
}

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = 9 * time.Second
)

// 初始化
func NewConnection(conn *websocket.Conn, connID uint32, apis ziface.ImsgHandle) *Connection {

	c := &Connection{
		Conn:       conn,
		ConnID:     connID,
		msgHandler: apis,
		isClosed:   false,
		msgChan:    make(chan []byte, 0),
		ExitChan:   make(chan bool, 1),
	}
	return c
}

func (c *Connection) Start() {
	log.Println("reader run...,connect_id:", c.ConnID)
	go c.StartReader()
	go c.StartWriter()
}

func (c *Connection) StartReader() {

	defer log.Println("connid = ", c.ConnID, " reader is exit ,addr is ")
	defer c.Stop()
	c.Conn.SetReadLimit(int64(utils.GlobalObject.MaxPackageSize))

	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(
		func(string) error {
			c.Conn.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})

	for {
		msgType, message, err := c.Conn.ReadMessage()
		_ = msgType
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
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
		c.msgHandler.DoMsgHandler(req)

	}
}

func (c *Connection) StartWriter() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		//c.Stop()
		if err := recover(); err != nil {
			log.Println(string(debug.Stack()))
		}
	}()

	for {
		select {

		case proto, ok := <-c.msgChan:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Conn.WriteMessage(websocket.BinaryMessage, proto)

		case <-c.ExitChan:
			log.Println("write 退出")
			return
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			c.Conn.WriteMessage(websocket.PingMessage, nil)

		}
	}
}

func (c *Connection) Stop() {
	if c.isClosed {
		return
	}
	c.isClosed = false
	c.Conn.Close()
	close(c.ExitChan)
	close(c.msgChan)
}

func (c *Connection) GetTCPConnection() *websocket.Conn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
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
