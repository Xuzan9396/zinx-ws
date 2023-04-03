package main

import (
	"flag"
	"github.com/Xuzan9396/zinx-ws/v5/znet"
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
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	p := znet.NewDataPack()
	by := []byte{'h', 'e', 'l', 'l', 'o'}

	resBytes, err := p.Pack(&znet.Message{
		Id:      1,
		DataLen: uint32(len(by)),
		Data:    by,
	})
	for {
		select {

		case <-ticker.C:
			sendMsg := resBytes
			err := c.WriteMessage(websocket.BinaryMessage, sendMsg)
			if err != nil {
				log.Println("write:", err)
				return
			}

			log.Println("写入成功:", string(by))
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
