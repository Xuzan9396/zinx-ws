package znet

import (
	"github.com/Xuzan9396/ws/v6/ziface"
	"github.com/gorilla/websocket"
	"log"
	"net"
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

	// 告知当前链接已经退出的 channel
	ExitChan chan bool
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = 9 * time.Second
	maxMessageSize = 10 * 1024
)

// 初始化
func NewConnection(conn *websocket.Conn, connID uint32, apis ziface.ImsgHandle) *Connection {

	c := &Connection{
		Conn:       conn,
		ConnID:     connID,
		msgHandler: apis,
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
	}
	return c
}

func (c *Connection) Start() {
	log.Println("reader run...,connect_id:", c.ConnID)
	defer log.Println("connid = ", c.ConnID, " reader is exit ,addr is ")
	defer c.Stop()
	c.Conn.SetReadLimit(maxMessageSize)
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

func (c *Connection) Stop() {
	if c.isClosed {
		return
	}
	c.isClosed = false
	c.Conn.Close()
	close(c.ExitChan)
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

func (c *Connection) Send(data []byte) error {
	return nil
}
