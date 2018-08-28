package impl

import (
	"errors"
	"sync"

	"github.com/gorilla/websocket"
)

type Connection struct {
	wsConn    *websocket.Conn
	outChan   chan []byte
	inChan    chan []byte
	closeChan chan struct{}
	isClose   bool
	lock      sync.Mutex
}

func NewConn(wsConn *websocket.Conn) (conn *Connection, err error) {
	conn = &Connection{
		wsConn:    wsConn,
		outChan:   make(chan []byte, 1024),
		inChan:    make(chan []byte, 1024),
		closeChan: make(chan struct{}),
	}

	go conn.readLoop()
	go conn.writeLoop()
	return
}

func (c *Connection) ReadMessage() (message []byte, err error) {
	select {
	case message = <-c.inChan:
	case <-c.closeChan:
		err = errors.New("connect is closed")
	}
	return
}

func (c *Connection) WriteMessage(message []byte) (err error) {
	select {
	case c.outChan <- message:
	case <-c.closeChan:
		err = errors.New("connect is closed")
	}
	return
}

func (c *Connection) Close() {
	c.wsConn.Close()

	c.lock.Lock()
	if !c.isClose {
		close(c.closeChan)
	}
	c.lock.Unlock()
}

func (c *Connection) readLoop() {
	var (
		data []byte
		err  error
	)
	for {

		if _, data, err = c.wsConn.ReadMessage(); err != nil {
			goto ERR
		}
		select {
		case c.inChan <- data:
		case <-c.closeChan:
			goto ERR
		}
	}
ERR:
	c.Close()
}

func (c *Connection) writeLoop() {
	var (
		data []byte
		err  error
	)
	defer c.Close()
	for {
		select {
		case data = <-c.outChan:
		case <-c.closeChan:
			return
		}
		if err = c.wsConn.WriteMessage(websocket.TextMessage, data); err != nil {
			return
		}
	}
}
