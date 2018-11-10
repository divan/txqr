package protocol

import (
	"fmt"
	"io"
	"io/ioutil"
)

// Encoder represents protocol encoder.
type Encoder struct {
	chunkLen int
}

// NewEncoder creates and inits a new encoder for the given chunk length.
func NewEncoder(n int) *Encoder {
	return &Encoder{
		chunkLen: n,
	}
}

// EncodeReader encodes data from reader and splits it into chunks to be
// futher converted to QR code frames.
func (e *Encoder) EncodeReader(r io.Reader) ([]string, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read: %v", err)
	}

	return e.Encode(string(data))
}

// Encode encodes data from reader and splits it into chunks to be
// futher converted to QR code frames.
func (e *Encoder) Encode(str string) ([]string, error) {
	if len(str) < e.chunkLen {
		return []string{str}, nil
	}

	numChunks := len(str)/e.chunkLen + 1

	// TODO(divan): use sync.Pool as this probably will be used many times
	ret := make([]string, numChunks)
	count := 0
	for start := 0; start < len(str); start += e.chunkLen {
		end := start + e.chunkLen
		if end > len(str) {
			end = len(str)
		}

		ret[count] = e.frame(start, len(str), str[start:end])
		count++
	}
	return ret, nil
}

func (e *Encoder) frame(count, total int, str string) string {
	return fmt.Sprintf("%x/%x|%s", count, total, str)
}
