package txqr // package should have this name to work properly with gomobile

import (
	"github.com/divan/txqr/protocol"
)

// Decoder implements txqr wrapper around protocol decoder.
type Decoder struct {
	*protocol.Decoder
}

// NewDecoder creats new txqr decoder.
func NewDecoder() *Decoder {
	return &Decoder{
		Decoder: protocol.NewDecoder(),
	}
}
