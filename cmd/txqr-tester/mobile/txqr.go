package txqrtester

// Connector represents a connection to the txqr-tester app
// via websockets.
type Connector struct {
}

// Connect attempts to establish connection to the WS server.
func (c *Connector) Connect(address string) error {
	return nil
}

// StartNext notifies txqr-tester that app is ready to scan next animated QR.
func (c *Connector) StartNext() error {
	return nil
}

// SendResult sends scanning result in milliseconds.
func (c *Connector) SendResult(duration int) error {
	return nil
}
