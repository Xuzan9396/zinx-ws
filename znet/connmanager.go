package znet

import (
	"errors"
	"github.com/Xuzan9396/ws/ziface"
	"sync"
)

type Connmanager struct {
	connections map[uint32]ziface.IConnection // 管理连接集合
	sync.RWMutex
}

// 穿件当前链接方法

func NewConnManager() *Connmanager {
	return &Connmanager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

func (c *Connmanager) Add(conn ziface.IConnection) {
	c.Lock()
	defer c.Unlock()
	c.connections[conn.GetConnID()] = conn
}

// 删除链接
func (c *Connmanager) Remove(conn ziface.IConnection) {
	c.Lock()
	defer c.Unlock()
	delete(c.connections, conn.GetConnID())
}

// 根据connID 获取链接
func (c *Connmanager) Get(connID uint32) (ziface.IConnection, error) {
	c.RLock()
	defer c.RUnlock()
	connt, ok := c.connections[connID]
	if !ok {
		return nil, errors.New("没有该链接")
	}
	return connt, nil
}

// 得到当前链接总数
func (c *Connmanager) Len() int {
	c.RLock()
	defer c.RUnlock()
	return len(c.connections)
}

// 清楚终止所有连接
func (c *Connmanager) ClearConn() {
	c.Lock()
	defer c.Unlock()
	for id, connection := range c.connections {
		connection.Stop()
		delete(c.connections, id)
	}

}
