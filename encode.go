package txqr

import (
	"fmt"
	"math/rand"

	fountain "github.com/google/gofountain"
)

// Encoder represents protocol encoder.
type Encoder struct {
	chunkLen         int
	redundancyFactor float64
}

// NewEncoder creates and inits a new encoder for the given chunk length.
func NewEncoder(n int) *Encoder {
	return &Encoder{
		chunkLen:         n,
		redundancyFactor: 2.0,
	}
}

// Encode encodes data from reader and splits it into chunks to be
// futher converted to QR code frames.
func (e *Encoder) Encode(str string) ([]string, error) {
	if len(str) < e.chunkLen {
		return []string{e.frame(0, len(str), []byte(str))}, nil
	}

	numChunks := numberOfChunks(len(str), e.chunkLen)
	codec := fountain.NewLubyCodec(numChunks, rand.New(fountain.NewMersenneTwister(200)), solitonDistribution(numChunks))

	var msg = []byte(str) // copy of str, as EncodeLTBlock is destructive to msg
	idsToEncode := ids(int(float64(numChunks) * e.redundancyFactor))
	lubyBlocks := fountain.EncodeLTBlocks(msg, idsToEncode, codec)

	// TODO(divan): use sync.Pool as this probably will be used many times
	ret := make([]string, len(lubyBlocks))
	for i, block := range lubyBlocks {
		ret[i] = e.frame(block.BlockCode, len(str), block.Data)
	}
	return ret, nil
}

// SetRedundancyFactor changes the value of redundancy factor.
func (e *Encoder) SetRedundancyFactor(rf float64) {
	e.redundancyFactor = rf
}

func (e *Encoder) frame(blockCode int64, total int, data []byte) string {
	return fmt.Sprintf("%d/%d/%d|%s", blockCode, e.chunkLen, total, string(data))
}

func numberOfChunks(length, chunkLen int) int {
	n := length / chunkLen
	if length%chunkLen > 0 {
		n++
	}
	return n
}
