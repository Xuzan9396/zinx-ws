package main

import (
	"flag"
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

	for {
		select {

		case <-ticker.C:
			sendMsg := `{"msgType":"ping"}`
			err := c.WriteMessage(websocket.BinaryMessage, []byte(sendMsg))
			if err != nil {
				log.Println("write:", err)
				return
			}
			log.Println("写入成功:", sendMsg)

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
