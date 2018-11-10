package protocol

import (
	"fmt"
	"sort"
)

// Decoder represents protocol decode.
type Decoder struct {
	buffer   []byte
	complete bool
	total    int
	frames   []frameInfo
}

// frameInfo represents the information about read frames.
// As frames can change size dynamically, we keep size info as well.
type frameInfo struct {
	offset, size int
}

// NewDecoder creats and inits a new decoder.
func NewDecoder() *Decoder {
	return &Decoder{
		buffer: []byte{},
	}
}

// NewDecoderSize creats and inits a new decoder for the known size.
// Note, it doesn't limit the size of the input, but optimizes memory allocation.
func NewDecoderSize(size int) *Decoder {
	return &Decoder{
		buffer: make([]byte, size),
	}
}

// DecodeChunk takes a single chunk of data and decodes it.
func (d *Decoder) DecodeChunk(data string) error {
	if data == "" || len(data) < 4 {
		return fmt.Errorf("invalid frame: \"%s\"", data)
	}

	var (
		offset, total int
		payload       []byte
	)
	_, err := fmt.Sscanf(data, "%d/%d|%s", &offset, &total, &payload)
	if err != nil {
		return fmt.Errorf("invalid frame: %v (%s)", err, data)
	}

	// allocate enough memory at first total read
	if total > d.total {
		d.buffer = make([]byte, total)
		d.total = total
	}

	size := len(payload)
	// TODO(divan): optmize memory allocation
	d.frames = append(d.frames, frameInfo{offset, size})

	copy(d.buffer[offset:offset+size], payload)

	// run the integrity check
	d.complete = d.isCompleted()

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

// IsCompleted reports whether the read was completed successfully or not.
func (d *Decoder) IsCompleted() bool {
	return d.complete
}

// isCompleted checks if all frames has been read.
// FIXME(divan): this approach might give false negatives in extreme cases, like
// many dynamic changes of chunk sizes.
func (d *Decoder) isCompleted() bool {
	sort.Slice(d.frames, func(i, j int) bool {
		return d.frames[i].offset < d.frames[j].offset
	})

	var cur int
	for _, frame := range d.frames {
		// we found the gap, next frame starts farther then current position
		if frame.offset > cur {
			return false
		}

		cur += frame.size
	}

	return cur == d.total
}
