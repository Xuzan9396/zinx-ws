package znet

import (
	"fmt"
	wsconfig "github.com/Xuzan9396/ws/config"
	"github.com/Xuzan9396/ws/ziface"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
	"runtime/debug"
	"sync/atomic"
	"time"
)

type Server struct {
	// 服务名称
	Name string
	// tcp4
	IPVersion string
	// ip地址
	IP string
	//端口
	Port int
	// 连接ID递增
	ConnIdInt uint32
	// 消息管理器
	msgHandler ziface.ImsgHandle
	// 链接管理模块
	ConnMgr ziface.IConnManager
	// 启动构造函数
	OnConnStart func(conn ziface.IConnection)
	//销毁构造函数
	OnConnStop func(conn ziface.IConnection)
}

var wsUpgrader = websocket.Upgrader{
	//	// 允许所有CORS跨域请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024 * 4,
	WriteBufferSize: 1024 * 4,
}

func NewServer() ziface.IServer {
	s := &Server{
		Name:       wsconfig.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         wsconfig.GlobalObject.Host,
		Port:       wsconfig.GlobalObject.TcpPort,
		msgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}
	return s
}

func (s *Server) Server() {
	log.Printf("version:%s,IP:%s,Port:%d", s.Name, s.IP, s.Port)
	s.Start()
}

func (s *Server) Start() {
	addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		log.Panicln(err)
		return
	}
	listerner, err := net.ListenTCP(s.IPVersion, addr)
	if err != nil {
		log.Panicln(err)
		return
	}

	r := s.routes()
	srv := &http.Server{Handler: r, ReadTimeout: 60 * time.Second}
	log.Println("success ws!!")
	err = srv.Serve(listerner)
	if err != nil {
		log.Panicln(err)
	}

}

// 设置链接地址
func (s *Server) routes() (r *mux.Router) {
	r = mux.NewRouter()
	r.HandleFunc("/", s.wSHandler).Methods("GET")
	return
}

func (s *Server) wSHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(string(debug.Stack()))
		}
	}()
	if r.Method != "GET" {
		return
	}

	if s.ConnMgr.Len() >= wsconfig.GlobalObject.MaxConn {
		log.Println("超过最大连接数:", s.ConnMgr.Len())
		return
	}

	socket, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	connIdInt := atomic.AddUint32(&s.ConnIdInt, 1)
	c := NewConnection(s, socket, connIdInt)
	log.Printf("连接成功,客户端IP:%s,conn_id:%d", c.RemoteAddr().String(), connIdInt)
	c.Start() // 启动读数据和写数据

}

func (s *Server) Stop() {
	s.ConnMgr.ClearConn()
}

func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.msgHandler.AddRouter(msgId, router)
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}
func (s *Server) GetMsgHandler() ziface.ImsgHandle {
	return s.msgHandler
}

//OnConnStart func(conn ziface.IConnection)
func (s *Server) SetOnConnStart(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

func (s *Server) SetOnConnStop(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

func (s *Server) GetOnConnStart() (hookFunc func(conn ziface.IConnection)) {
	return s.OnConnStart
}

func (s *Server) GetOnConnStop() (hookFunc func(conn ziface.IConnection)) {
	return s.OnConnStop
}
