package protocol

import (
	"fmt"
	"strings"
	"time"

	"github.com/pyk/byten"
)

// Decoder represents protocol decode.
type Decoder struct {
	buffer   []byte
	complete bool
	total    int
	frames   []frameInfo
	cache    map[string]struct{}

	progress int
	speed    int // avg reading speed
	start    time.Time

	lastChunk    time.Time // last chunk decode request
	readInterval time.Duration
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
		cache:  make(map[string]struct{}),
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
	if d.IsCompleted() {
		return nil
	}
	if !d.lastChunk.IsZero() {
		d.readInterval = time.Now().Sub(d.lastChunk)
	}
	d.lastChunk = time.Now()

	if data == "" || len(data) < 4 {
		return fmt.Errorf("invalid frame: \"%s\"", data)
	}

	idx := strings.IndexByte(data, '|')
	if idx == -1 {
		return fmt.Errorf("invalid frame: \"%s\"", data)
	}

	if d.start.IsZero() {
		d.start = time.Now()
	}

	header := data[:idx]
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

	payload := data[idx+1:]
	size := len(payload)
	// TODO(divan): optmize memory allocation
	d.frames = append(d.frames, frameInfo{offset: offset, size: size})

	copy(d.buffer[offset:offset+size], payload)

	d.updateProgress()

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

// Speed returns avg reading speed.
func (d *Decoder) Speed() string {
	return fmt.Sprintf("%s/s", byten.Size(int64(d.speed)))
}

// ReadInterval returns the latest read interval in ms.
func (d *Decoder) ReadInterval() int64 {
	return int64(d.readInterval / time.Millisecond)
}

// Progress returns reading progress in percentage.
func (d *Decoder) Progress() int {
	return d.progress
}

// TotalTime returns the total scan duration in human readable form - from first to last read chunk.
func (d *Decoder) TotalTime() string {
	dur := time.Since(d.start)
	return formatDuration(dur)
}

// formatDuration converts "12.232312313s" to "12.2s"
func formatDuration(d time.Duration) string {
	if d > time.Second {
		d = d - d%(100*time.Millisecond)
	}
	return d.String()
}

// TotalSize returns the data size in human readable form.
func (d *Decoder) TotalSize() string {
	return byten.Size(int64(len(d.buffer)))
}

// IsCompleted reports whether the read was completed successfully or not.
func (d *Decoder) IsCompleted() bool {
	return d.complete
}

// updateProgress updates progress and complete state of reading.
// FIXME(divan): this approach might give false negatives in extreme cases, like
// many dynamic changes of chunk sizes.
func (d *Decoder) updateProgress() {
	var cur int
	for _, frame := range d.frames {
		cur += frame.size
	}

	d.speed = cur * int(time.Second) / int(time.Since(d.start))
	d.progress = 100 * cur / d.total
	d.complete = cur == d.total
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

// TotalTimeMs returns the total scan duration in milliseconds.
func (d *Decoder) TotalTimeMs() int64 {
	dur := time.Since(d.start)
	return int64(dur / time.Millisecond)
}

// Reset resets decoder, preparing it for the next run.
func (d *Decoder) Reset() {
	d.buffer = []byte{}
	d.complete = false
	d.total = 0
	d.frames = []frameInfo{}
	d.cache = map[string]struct{}{}

	d.progress = 0
	d.speed = 0
	d.start = time.Time{}
	d.lastChunk = time.Time{}
	d.readInterval = 0
}
