package txqrtester

import (
	"encoding/json"
	"fmt"

	"github.com/divan/txqr/cmd/txqr-tester/ws"
	"github.com/divan/txqr/mobile"
	"github.com/gorilla/websocket"
)

// Connector represents a connection to the txqr-tester app
// via websockets.
type Connector struct {
	conn *websocket.Conn
	*txqr.Decoder
}

func NewConnector() *Connector {
	return &Connector{
		Decoder: txqr.NewDecoder(),
	}
}

// Connect attempts to establish connection to the WS server.
func (c *Connector) Connect(address string) error {
	conn, _, err := websocket.DefaultDialer.Dial(address, nil)
	if err != nil {
		return fmt.Errorf("dial: %v", err)
	}
	c.conn = conn
	return c.connect()
}

// StartNext notifies txqr-tester that app is ready to scan next animated QR.
func (c *Connector) StartNext() error {
	return c.sendCommand(ws.UIRequest{
		Cmd: ws.CmdStartNext,
	})
}

// SendResult sends scanning result in milliseconds.
func (c *Connector) SendResult(duration int64) error {
	return c.sendCommand(ws.UIRequest{
		Cmd:      ws.CmdResult,
		Duration: duration,
	})
}

// Close closes underlying WebSocket connection.
func (c *Connector) Close() error {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

func (c *Connector) connect() error {
	return c.sendCommand(ws.UIRequest{
		Cmd:    ws.CmdConnect,
		Device: "N/A",
	})
}

func (c *Connector) sendCommand(req ws.UIRequest) error {
	data, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("json marshal: %v", err)
	}
	return c.conn.WriteMessage(websocket.TextMessage, data)
}
