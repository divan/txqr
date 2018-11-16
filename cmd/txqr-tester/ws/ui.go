package ws

type (
	// UIResponse is a WebSocket UI response definition.
	UIResponse struct {
		Type MsgType `json:"type"`
	}

	// UIRequest is a WebSocket UI request definition.
	UIRequest struct {
		Cmd    UICommand `json:"cmd"`
		Device string    `json:"device,omitempty"`
	}

	// UICommand defines client commands.
	UICommand string
	MsgType   string
)

var (
	CmdConnect UICommand = "connect"

	TypeAck MsgType = "ack"
)
