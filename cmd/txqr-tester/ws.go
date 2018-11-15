package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// WSServer implements WebSocket server for receiving
// QR scan state and results.
type WSServer struct {
	upgrader  websocket.Upgrader
	connected bool
}

func NewWSServer() *WSServer {
	ws := &WSServer{
		upgrader: websocket.Upgrader{},
	}
	return ws
}

type WSResponse struct {
	Type MsgType `json:"type"`
}

type WSRequest struct {
	Cmd WSCommand `json:"cmd"`
}

type MsgType string
type WSCommand string

// WebSocket response types
const (
	RespPositions MsgType = "positions"
	RespGraph     MsgType = "graph"
	RespStats     MsgType = "stats"
)

// WebSocket commands
const (
	CmdHello  WSCommand = "hello"
	CmdBye    WSCommand = "bye"
	CmdResult WSCommand = "result"
)

func (ws *WSServer) Handle(w http.ResponseWriter, r *http.Request) {
	c, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer c.Close()

	ws.connected = true

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", mt, err)
			break
		}
		ws.processRequest(c, mt, message)
	}
}

func (ws *WSServer) processRequest(c *websocket.Conn, mtype int, data []byte) {
	var cmd WSRequest
	err := json.Unmarshal(data, &cmd)
	if err != nil {
		log.Println("[ERROR] invalid command:", err)
		return
	}

	switch cmd.Cmd {
	case CmdHello:
		fmt.Println("Got hello")
	case CmdBye:
		fmt.Println("Got bye")
	case CmdResult:
		fmt.Println("Got Result")
	}
}

func (ws *WSServer) sendMsg(c *websocket.Conn, msg *WSResponse) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Println("write:", err)
		return
	}

	err = c.WriteMessage(1, data)
	if err != nil {
		log.Println("write:", err)
		return
	}
}
