package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// WSBridge implements WebSocket bridge for connecting
// two clients directly to each other.
type WSBridge struct {
	upgrader websocket.Upgrader

	first, second *websocket.Conn
}

func NewWSBridge() *WSBridge {
	ws := &WSBridge{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
	return ws
}

func (ws *WSBridge) Handle(w http.ResponseWriter, r *http.Request) {
	log.Println("[DEBUG] Got new connection from", r.Host)
	conn, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	if ws.first == nil {
		ws.first = conn
		ws.handleFirstClient()
	} else if ws.second == nil {
		ws.second = conn
		ws.handleSecondClient()
	} else {
		log.Println("[ERROR] Already have two clients connected, dropping new connection from", r.Host)
		return
	}

}

func (ws *WSBridge) handleFirstClient() {
	for {
		mt, message, err := ws.first.ReadMessage()
		if err != nil {
			log.Println("[ERROR] Read first:", err)
			break
		}
		if ws.second == nil {
			// TODO(divan): send not ready?
			log.Println("Second client is not ready, skipping...")
			continue
		}
		err = ws.second.WriteMessage(mt, message)
		if err != nil {
			log.Println("[ERROR] Write second:", err)
			return
		}
	}
}

func (ws *WSBridge) handleSecondClient() {
	for {
		mt, message, err := ws.second.ReadMessage()
		if err != nil {
			log.Println("[ERROR] Read second:", err)
			break
		}
		if ws.first == nil {
			// TODO(divan): send not ready?
			log.Println("First client is not ready, skipping...")
			continue
		}
		err = ws.first.WriteMessage(mt, message)
		if err != nil {
			log.Println("[ERROR] Write first:", err)
			return
		}
	}
}
