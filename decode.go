package txqr

import (
	"fmt"
	"strings"
)

// Decoder represents protocol decode.
type Decoder struct {
	buffer      []byte
	read, total int
	frames      []frameInfo
	cache       map[string]struct{}
}

// frameInfo represents the information about read frames.
// As frames can change size dynamically, we keep size info as well.
type frameInfo struct {
	offset, size int
}

// NewDecoder creates and inits a new decoder.
func NewDecoder() *Decoder {
	return &Decoder{
		buffer: []byte{},
		cache:  make(map[string]struct{}),
	}
}

// NewDecoderSize creates and inits a new decoder for the known size.
// Note, it doesn't limit the size of the input, but optimizes memory allocation.
func NewDecoderSize(size int) *Decoder {
	return &Decoder{
		buffer: make([]byte, size),
	}
}

// Decode takes a single chunk of data and decodes it.
// Chunk expected to be validated (see Validate) before.
func (d *Decoder) Decode(chunk string) error {
	idx := strings.IndexByte(chunk, '|') // expected to be validated before
	header := chunk[:idx]

	// continuous QR reading often sends the same chunk in a row, skip it
	if d.isCached(header) {
		return nil
	}

	var offset, total int
	_, err := fmt.Sscanf(header, "%d/%d", &offset, &total)
	if err != nil {
		return fmt.Errorf("invalid header: %v (%s)", err, header)
	}

	// allocate enough memory at first total read
	if d.total == 0 {
		d.buffer = make([]byte, total)
		d.total = total
	}

	if total > d.total {
		return fmt.Errorf("total changed during sequence, aborting")
	}

	payload := chunk[idx+1:]
	size := len(payload)
	// TODO(divan): optmize memory allocation
	d.frames = append(d.frames, frameInfo{offset: offset, size: size})

	copy(d.buffer[offset:offset+size], payload)

	d.updateCompleted()

	return nil
}

// Validate checks if a given chunk of data is a valid txqr protocol packet.
func (d *Decoder) Validate(chunk string) error {
	if chunk == "" || len(chunk) < 4 {
		return fmt.Errorf("invalid frame: \"%s\"", chunk)
	}

	idx := strings.IndexByte(chunk, '|')
	if idx == -1 {
		return fmt.Errorf("invalid frame: \"%s\"", chunk)
	}

	return nil
}

// Data returns decoded data.
func (d *Decoder) Data() string {
	return string(d.buffer)
}

// DataBytes returns decoded data as a byte slice.
func (d *Decoder) DataBytes() []byte {
	return d.buffer
}

// Length returns length of the decoded data.
func (d *Decoder) Length() int {
	return len(d.buffer)
}

// Read returns amount of currently read bytes.
func (d *Decoder) Read() int {
	return d.read
}

// Total returns total amount of data.
func (d *Decoder) Total() int {
	return d.total
}

// IsCompleted reports whether the read was completed successfully or not.
func (d *Decoder) IsCompleted() bool {
	return d.total > 0 && d.read >= d.total
}

// isCached takes the header of chunk data and see if it's been cached.
// If not, it caches it.
func (d *Decoder) isCached(header string) bool {
	if _, ok := d.cache[header]; ok {
		return true
	}

	// cache it
	d.cache[header] = struct{}{}
	return false
}

// Reset resets decoder, preparing it for the next run.
func (d *Decoder) Reset() {
	d.buffer = []byte{}
	d.read, d.total = 0, 0
	d.frames = []frameInfo{}
	d.cache = map[string]struct{}{}
}

// TODO(divan): this will now work if frame size is dynamic. Rewrite it
// to support it.
func (d *Decoder) updateCompleted() {
	var cur int
	for _, frame := range d.frames {
		cur += frame.size
	}

	d.read = cur
}
