package main

import (
	"fmt"
	"imooc_learning/go-websocket/impl"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// upgrader.Upgrade()
	var (
		conn   *impl.Connection
		wsConn *websocket.Conn
		err    error
		data   []byte
	)

	if wsConn, err = upgrader.Upgrade(w, r, nil); err != nil {
		fmt.Println(err)
		return
	}

	if conn, err = impl.NewConn(wsConn); err != nil {
		return
	}

	go func() {
		t := time.NewTicker(time.Second)
		for range t.C {
			conn.WriteMessage([]byte("ping pong"))
		}
	}()
	for {

		if data, err = conn.ReadMessage(); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(data))
		if err = conn.WriteMessage(data); err != nil {
			return
		}
	}
}

func main() {
	http.HandleFunc("/ws", wsHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("ListenAndServe Error", err)
	}
}
