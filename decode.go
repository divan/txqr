package txqr

import (
	"fmt"
	"math/rand"
	"strings"

	fountain "github.com/google/gofountain"
)

// Decoder represents protocol decode.
type Decoder struct {
	chunkLen  int
	codec     fountain.Codec
	fd        fountain.Decoder
	completed bool
	total     int
	cache     map[string]struct{}
}

// NewDecoder creates and inits a new decoder.
func NewDecoder() *Decoder {
	return &Decoder{
		cache: make(map[string]struct{}),
	}
}

// NewDecoderSize creates and inits a new decoder for the known size.
func NewDecoderSize(size, chunkLen int) *Decoder {
	numChunks := numberOfChunks(size, chunkLen)
	codec := fountain.NewLubyCodec(numChunks, rand.New(fountain.NewMersenneTwister(200)), solitonDistribution(numChunks))
	return &Decoder{
		codec:    codec,
		fd:       codec.NewDecoder(size),
		total:    size,
		chunkLen: chunkLen,
		cache:    make(map[string]struct{}),
	}
}

// Decode takes a single chunk of data and decodes it.
// Chunk expected to be validated (see Validate) before.
func (d *Decoder) Decode(chunk string) error {
	idx := strings.IndexByte(chunk, '|') // expected to be validated before
	if idx == -1 {
		return fmt.Errorf("invalid frame: \"%s\"", chunk)
	}

	header := chunk[:idx]
	// continuous QR reading often sends the same chunk in a row, skip it
	if d.isCached(header) {
		return nil
	}

	var (
		blockCode       int64
		chunkLen, total int
	)
	_, err := fmt.Sscanf(header, "%d/%d/%d", &blockCode, &chunkLen, &total)
	if err != nil {
		return fmt.Errorf("invalid header: %v (%s)", err, header)
	}

	payload := chunk[idx+1:]
	lubyBlock := fountain.LTBlock{
		BlockCode: blockCode,
		Data:      []byte(payload),
	}

	if d.fd == nil {
		d.total = total
		d.chunkLen = chunkLen
		numChunks := numberOfChunks(d.total, d.chunkLen)
		d.codec = fountain.NewLubyCodec(numChunks, rand.New(fountain.NewMersenneTwister(200)), solitonDistribution(numChunks))
		d.fd = d.codec.NewDecoder(total)
	}
	d.completed = d.fd.AddBlocks([]fountain.LTBlock{lubyBlock})

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
	return string(d.DataBytes())
}

// DataBytes returns decoded data as a byte slice.
func (d *Decoder) DataBytes() []byte {
	if d.fd == nil {
		return []byte{}
	}

	if !d.completed {
		return []byte{}
	}
	return d.fd.Decode()
}

// Length returns length of the decoded data.
// TODO: remove
func (d *Decoder) Length() int {
	return 0
}

// Read returns amount of currently read bytes.
// TODO: remove
func (d *Decoder) Read() int {
	return 0
}

// Total returns total amount of data.
func (d *Decoder) Total() int {
	return d.total
}

// IsCompleted reports whether the read was completed successfully or not.
func (d *Decoder) IsCompleted() bool {
	return d.completed
}

// Reset resets decoder, preparing it for the next run.
func (d *Decoder) Reset() {
	d.fd = nil
	d.completed = false
	d.chunkLen = 0
	d.total = 0
	d.cache = map[string]struct{}{}
	d.codec = nil
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
