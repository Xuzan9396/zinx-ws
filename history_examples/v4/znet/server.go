package znet

import (
	"fmt"
	"github.com/Xuzan9396/zinx-ws/v4/ziface"
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
	Name      string
	IPVersion string
	IP        string
	Port      int
	ConnIdInt uint32
	// 对应router
	Router ziface.IRouter
}

var wsUpgrader = websocket.Upgrader{
	//	// 允许所有CORS跨域请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024 * 4,
	WriteBufferSize: 1024 * 4,
}

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "127.0.0.1",
		Port:      8999,
	}
	return s
}

func (s *Server) Server() {
	s.Start()
}

func (s *Server) Start() {
	log.Printf("version:%s,IP:%s,Port:%d", s.Name, s.IP, s.Port)
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

	socket, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	connIdInt := atomic.AddUint32(&s.ConnIdInt, 1)
	c := NewConnection(socket, connIdInt, s.Router)

	log.Printf("连接成功,客户端IP:%s,conn_id:%d", c.RemoteAddr().String(), connIdInt)
	c.Start() // 读取数据

}

func (s *Server) Stop() {

}

func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
}

//func CallBackToClient(conn *websocket.Conn, data []byte, cnt int) error {
//
//	log.Printf("msgType:%d,msg:%s \n", cnt, data)
//	return nil
//}
